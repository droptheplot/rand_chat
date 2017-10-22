package models

import (
	"os"
	"testing"

	"github.com/droptheplot/rand_chat/env"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	db = env.Init()
	env.Reset(db)

	os.Exit(m.Run())
}
