package core

import (
	"fmt"
	"time"

	config "github.com/karthikic/techblogs/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Blogs struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Title     string
	Image     string
	Company   string
	Link      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func SetupDB() *gorm.DB {
	db_configs := config.GetDbConfigs()

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		db_configs.User, db_configs.Password, db_configs.Host, db_configs.Port, db_configs.Database)

	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to open DB connection")
	}

	fmt.Println("Connection to DB Opened:", db)

	// declare migrations
	db.AutoMigrate(&Blogs{})
	return db
}
