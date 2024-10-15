package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint    `gorm:"primaryKey"`
	Email        string  `gorm:"size:255"`
	UserSettings YamlMap `gorm:"type:text"`
}

var db *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Microsecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,             // Disable color
		},
	)
	database, err := gorm.Open(mysql.Open("root:devdba_sys@tcp(127.0.0.1:3307)/gorm_yaml_serializer?charset=utf8&parseTime=true"), &gorm.Config{Logger: newLogger})

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database!")
	}

	db = database
}

func dbMigrate() {
	db.AutoMigrate(&User{})
}

func main() {
	connectDatabase()
	dbMigrate()

	user := User{
		Email: "abc@codeheim.io",
		UserSettings: YamlMap{
			"theme":         "dark",
			"notifications": true,
		},
	}

	db.Create(&user)

	var fetchedUser User
	db.First(&fetchedUser, user.ID)

	fmt.Println(fetchedUser)
}
