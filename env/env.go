package env

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("postgres", Config.Database)

	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
}

func Migrate() {
	driver, _ := postgres.WithInstance(DB.DB(), &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	migrations.Up()
}
