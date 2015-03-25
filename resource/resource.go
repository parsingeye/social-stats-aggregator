package resource

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/guilex/social-stats-aggregator/api"
	"github.com/guilex/social-stats-aggregator/model"
	"github.com/guilex/social-stats-aggregator/retriever"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type IndexStatForm struct {
	Id       []string `form:"id[]" binding:"required"`
	Provider string   `form:"provider" binding:"required"`
}

type StoreStatForm struct {
	Url      string `form:"url" binding:"required"`
	Provider string `form:"provider" binding:"required"`
}

type StatResource struct {
	Db gorm.DB
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

	provider, _ := retriever.Make(s.Provider)

	likes := provider.GetCount(s.Url)

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

func (resource *StatResource) Index(c *gin.Context) {

	resource.Db.LogMode(true)

	c.Request.ParseForm()

	var form IndexStatForm

	if c.BindWith(&form, binding.Form) {

		var stats struct {
			Data []model.Stat `json:"data"`
		}

		resource.Db.Where("provider = ? AND id IN (?)", form.Provider, form.Id).Find(&stats.Data)

		if len(stats.Data) == 0 {

			c.JSON(http.StatusNotFound, api.NewError("404 Not Found"))

		} else {

			c.JSON(http.StatusOK, stats)

		}

	} else {

		c.JSON(http.StatusBadRequest, api.NewError("400 Bad Request"))

	}
}

func (resource *StatResource) Show(c *gin.Context) {

	idString := c.Params.ByName("id")

	id, _ := strconv.Atoi(idString)

	fmt.Println("\v", id)

	var stat model.Stat

	if resource.Db.Where("id = ?", id).Find(&stat).RecordNotFound() {
		c.JSON(http.StatusNotFound, api.NewError("404 Not Found"))
	} else {
		c.JSON(http.StatusOK, stat)
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

		provider, _ := retriever.Make(form.Provider)

		count := provider.GetCount(form.Url)

		var stat model.Stat
		stat.Url = form.Url
		stat.Urlhash = urlHash
		stat.Count = count
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
