package service

import (
	"../model"
	"../resource"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

type Config struct {
	ServiceHost string
}

type StatService struct {
}

func (s *StatService) getDb(config Config) (gorm.DB, error) {
	return gorm.Open("sqlite3", "/Users/dragan/Sites/gotest/gorm.db")
}

func (s *StatService) Migrate(config Config) error {
	db, err := s.getDb(config)
	if err != nil {
		return err
	}

	db.AutoMigrate(&model.Stat{})
	return nil
}

func (s *StatService) Update(config Config, interval int) error {

	db, err := s.getDb(config)

	if err != nil {
		return err
	}

	statResource := &resource.StatResource{}

	statResource.Db = db

	statResource.BatchUpdate(interval)

	return nil
}

func (s *StatService) Run(config Config) error {

	db, err := s.getDb(config)

	if err != nil {
		return err
	}

	statResource := &resource.StatResource{}

	statResource.Db = db

	r := gin.Default()

	r.POST("/stats", statResource.Store)

	r.GET("/stats", statResource.Get)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(config.ServiceHost)

	return nil
}
