package configs

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB = connect()

func connect() *gorm.DB {
	HOST := "localhost"
	PORT := 5432
	DBNAME := "concurrente"
	USER := "postgres"
	PASS := "admin123"
	SSLMODE := "disable"

	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s",
		HOST, USER, PASS, DBNAME, PORT, SSLMODE)

	Postgres, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("No connection to the Database")
	}

	log.Println("Connected to the Database")

	return Postgres
}
