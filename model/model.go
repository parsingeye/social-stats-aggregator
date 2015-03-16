package model

import (
	"time"
)

type Stat struct {
	ID        int    `sql:"AUTO_INCREMENT"`
	Url       string `sql:"size:255" form:"url" binding:"required"`
	Urlhash   string `sql:"size:32"`
	Provider  string `sql:"size:128" form:"provider" binding:"required"`
	Count     int
	Interval  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
