package models

import (
	"time"

	"github.com/droptheplot/rand_chat/env"
)

type Chart struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func GetCharts() (charts []Chart) {
	env.DB.Find(&charts)

	return charts
}
