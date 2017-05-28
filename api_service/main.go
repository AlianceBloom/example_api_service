package api_service

import (
	"log"
	"github.com/aliancebloom/example_api_service/game"
)


func Run() {
	DB, err := initDatabase()
	if err != nil {
		log.Fatal("Can't connect to DB,", err)
	}
	defer DB.Close()

	game := game.InitializeGame(DB)
	log.Println(game)
	runServer(game)
	log.Println("Running...")
}
