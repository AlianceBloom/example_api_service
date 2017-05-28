package api_service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDatabase() (*gorm.DB, error) {
	return gorm.Open("sqlite3", "./tmp/game.db")
}
