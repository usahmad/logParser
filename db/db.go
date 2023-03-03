package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func InitDb() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		return nil
	}
	dbTable, exists := os.LookupEnv("DB_SCHEMA")
	if !exists {
		return nil
	}
	Db = connectDB(dbTable)
	return Db
}

func connectDB(dbName string) *gorm.DB {
	Username, exists := os.LookupEnv("DB_USERNAME")
	Password, exists := os.LookupEnv("DB_PASSWORD")
	Host, exists := os.LookupEnv("DB_HOST")
	Port, exists := os.LookupEnv("DB_PORT")
	if !exists {
		fmt.Printf("NO DATA IN ENV FILE")
		return nil
	}
	var err error
	dsn := Username + ":" + Password + "@tcp" + "(" + Host + ":" + Port + ")/" + dbName + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	return db
}
