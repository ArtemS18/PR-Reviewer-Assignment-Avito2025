package main

import (
	"reviewer-api/internal/app/config"
	"reviewer-api/internal/app/ds"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	config, err := config.New()
	if err != nil {
		panic("failed to load config")
	}
	db, err := gorm.Open(postgres.Open(config.GetDSN()))
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(
		&ds.PullRequest{},
		&ds.User{},
		&ds.Reviewer{},
		&ds.Team{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
