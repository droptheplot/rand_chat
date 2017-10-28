package app

import (
	"os"
	"testing"

	"github.com/droptheplot/rand_chat/env"
	"github.com/rs/zerolog"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var logger zerolog.Logger

func TestMain(m *testing.M) {
	db, logger = env.Init()
	env.Reset(db)

	os.Exit(m.Run())
}
