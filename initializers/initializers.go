package initializers

import (
	"fmt"
	"log"
	"os"

	"sportzone/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Loadenvariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Initdb() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("HOST"),
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("PORT"))
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("Error connecting to database")
	}
}

func Syncdatabase() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error creating database")
	}
}
