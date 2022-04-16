package models

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// Setup initializes the database instance
func Setup() {
	var err error

	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)

	db, err = gorm.Open(mysql.Open(url))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	if !db.Migrator().HasTable(&Auth{}) {
		db.AutoMigrate(&Auth{})
		if db.Migrator().HasTable(&Auth{}) {
			logging.Info("create config table success")
		} else {
			logging.Error("create config table fail")
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
