package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	//"os"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	//username := os.Getenv("db_user")
	//password := os.Getenv("db_pass")
	//dbName := os.Getenv("db_name")
	//dbHost := os.Getenv("db_host")


	dbUri := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		"localhost", "gdeon", "learning_go_db", "3228")
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{}, &Currency{})
}

func GetDB() *gorm.DB {
	return db
}