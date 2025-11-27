package main

import (
	"DeBroglieProject/internal/app/ds"
	"DeBroglieProject/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Particle{},
		&ds.RequestDeBroglieCalculation{},
		&ds.DeBroglieCalculation{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
