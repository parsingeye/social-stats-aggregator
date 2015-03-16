package resource

import (
	"../api"
	"../model"
	"../retriever"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type StatResource struct {
	Db gorm.DB
}

type GetStatForm struct {
	Url      string `form:"url" binding:"required"`
	Provider string `form:"provider"`
}

type StoreStatForm struct {
	Url      string `form:"url" binding:"required"`
	Provider string `form:"provider" binding:"required"`
}

func (resource *StatResource) BatchUpdate(interval int) {

	var stats struct {
		Data []model.Stat
	}

	fmt.Println("Start Batch UPDATE for interval")

	resource.Db.Where("interval = ?", interval).Find(&stats.Data)

	for _, s := range stats.Data {
		resource.update(s)
	}

}

func (resource *StatResource) update(s model.Stat) {

	if s.Provider == "facebook" {

		facebookRetriever := retriever.Facebook{}

		likes := facebookRetriever.GetLikes(s.Url)

		if s.Count != likes {
			s.Count = likes
			s.UpdatedAt = time.Now()
			s.Interval = resetInterval()
		} else {
			s.Interval = degradeInterval(s.Interval)
		}

		fmt.Println("\v", s)

		resource.Db.Save(&s)
	}
}

func (resource *StatResource) Get(c *gin.Context) {

	resource.Db.LogMode(true)

	c.Request.ParseForm()

	var form GetStatForm

	c.BindWith(&form, binding.Form)

	urlHash := getMD5Hash(form.Url)

	var stats struct {
		Data []model.Stat
	}

	fmt.Println("\v", form)

	query := resource.Db.Where("provider = ? and urlhash = ?", form.Provider, urlHash)

	if form.Provider == "" {

		query = resource.Db.Where("urlhash = ?", urlHash)

	}

	if query.Find(&stats.Data).RecordNotFound() {

		c.JSON(http.StatusNotFound, api.NewError("404 Not Found"))

	} else {

		c.JSON(http.StatusOK, stats)

	}
}

func (resource *StatResource) Store(c *gin.Context) {

	c.Request.ParseForm()

	var form StoreStatForm

	c.BindWith(&form, binding.MultipartForm)

	urlHash := getMD5Hash(form.Url)

	// check for existance
	count := 0

	resource.Db.Where("provider = ? and urlhash = ?", form.Provider, urlHash).Find(&model.Stat{}).Count(&count)

	if count > 0 {

		c.JSON(http.StatusConflict, api.NewError("Conflict"))
		return

	} else {

		facebookRetriever := retriever.Facebook{}

		likes := facebookRetriever.GetLikes(form.Url)

		var stat model.Stat
		stat.Url = form.Url
		stat.Urlhash = urlHash
		stat.Count = likes
		stat.Interval = 1
		stat.Provider = form.Provider
		stat.CreatedAt = time.Now()
		stat.UpdatedAt = time.Now()

		resource.Db.Create(&stat)

		c.JSON(http.StatusCreated, stat)
	}
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func resetInterval() int {
	return 1
}

func degradeInterval(current int) int {

	if newInterval := current + 1; newInterval <= 5 {
		return newInterval
	}

	return 5
}
