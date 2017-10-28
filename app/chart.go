package app

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Chart struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func GetCharts(db *gorm.DB) (charts []Chart) {
	db.Find(&charts)

	return charts
}
