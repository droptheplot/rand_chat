package env

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/rs/zerolog"
)

func Init() (*gorm.DB, zerolog.Logger) {
	db, err := gorm.Open("postgres", Config.Database)

	if err != nil {
		panic(err)
	}

	logger := zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return db, logger
}

func Migrate(db *gorm.DB) {
	driver, _ := postgres.WithInstance(db.DB(), &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance(Config.Migrations, "postgres", driver)
	migrations.Up()
}

func Drop(db *gorm.DB) {
	db.Exec(`DROP SCHEMA public CASCADE;`)
	db.Exec(`CREATE SCHEMA public;`)
}

func Reset(db *gorm.DB) {
	Drop(db)
	Migrate(db)
}
