package postgres

import (
	"log"
	"reviewer-api/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgers(dsn string, autoMigrate bool) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	if autoMigrate {
		err = db.AutoMigrate(
			&ds.PullRequest{},
			&ds.User{},
			&ds.Reviewer{},
			&ds.Team{},
		)
		if err != nil {
			panic("cant migrate db")
		}
		log.Println("success migrate db")
	}
	return &Postgres{db}, nil
}
