package api_service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
)

func initDatabase() (*gorm.DB, error) {
	tmpDir, _ := ioutil.TempDir("", "text_game_engine")
	return gorm.Open("sqlite3", tmpDir + "/game.sqlite")
}
