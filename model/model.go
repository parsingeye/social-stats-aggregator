package model

import (
	"time"
)

type Stat struct {
	ID        int       `sql:"AUTO_INCREMENT" json:"id"`
	Url       string    `sql:"size:255" form:"url" binding:"required" json:"url"`
	Urlhash   string    `sql:"size:32" json:"url_hash"`
	Provider  string    `sql:"size:128" form:"provider" binding:"required" json:"provider"`
	Count     int       `json:"count"`
	Interval  int       `json:"interval"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
}
