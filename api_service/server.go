package api_service

import (
	"github.com/gin-gonic/gin"
	"github.com/aliancebloom/example_api_service/game"
	"net/http"
	"log"
)

type ResultTournamentParams struct {
	TournamentId string `binding:"required"`
	PlayerId string `binding:"required"`
}

type NewPlayerParams struct {
	PlayerId string `binding:"required"`
	Points   int `binding:"required"`
}

type PlayerManipulationParams struct {
	PlayerId string `binding:"required"`
	Amount int `binding:"required"`
}

type NewTournamentParams struct {
	TournamentId string `binding:"required"`
	Deposit int `binding:"required"`
}

type JoinTournamentParams struct {
	TournamentId string `binding:"required"`
	PlayerId string `binding:"required"`
	BackerId []string
}

type PlayerParams struct {
	PlayerId string `binding:"required"`
}

var router *gin.Engine
var gameEngine *game.Game

func runServer(game *game.Game) {
	gameEngine = game

	router := gin.Default()

	router.POST("/newUser", newUserEndpoint)
	router.GET("/balance", balanceEndpoint)
	router.PUT("/take", takeEndpoint)
	router.PUT("/fund", fundEndpoint)
	router.POST("/announceTournament", announceTournamentEndpoint)
	router.PUT("/joinTournament", joinTournamentEndpoint)
	router.POST("/resultTournament", resultTournamentEndpoint)

	router.Run()
}

func balanceEndpoint(c *gin.Context) {
	var params PlayerParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	player := game.FindPlayer(params.PlayerId)
	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playerId": player.ID, "balance": player.Points})
}


func newUserEndpoint(c *gin.Context) {
	var params NewPlayerParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if _, err := gameEngine.AddPlayer(params.Points, params.PlayerId); err != nil {
		log.Println("Errors", err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}


func takeEndpoint(c *gin.Context) {
	var params PlayerManipulationParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	player := game.FindPlayer(params.PlayerId)
	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	_, err := player.TakePoints(params.Amount)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}


func fundEndpoint(c *gin.Context) {
	var params PlayerManipulationParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	player := game.FindPlayer(params.PlayerId)
	if player == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	_, err := player.FundPoints(params.Amount)
	if err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}


func announceTournamentEndpoint(c *gin.Context) {
	var params NewTournamentParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	_, err := gameEngine.NewTournament(params.TournamentId, params.Deposit)

	if err != nil {
		log.Println("Error:	", err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
	return

}


func joinTournamentEndpoint(c *gin.Context) {
	var params JoinTournamentParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	tournament := game.FindTournament(params.TournamentId)
	if tournament == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if err := tournament.JoinTournament(params.PlayerId, params.BackerId...); err != nil {
		log.Println("Error while joining tournament", err)
		c.JSON(http.StatusMethodNotAllowed, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}


func resultTournamentEndpoint(c *gin.Context) {
	var params ResultTournamentParams
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	tournament := game.FindTournament(params.TournamentId)
	if tournament == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	result, err := tournament.FinishTournament(params.PlayerId)
	if err != nil {
		log.Println("Error while finishing tournament", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	response := map[string]interface{}{"playerId": params.PlayerId, "prize": result.Amount }
	c.JSON(http.StatusCreated, gin.H{ "tournamentId": tournament.ID, "winners": []interface{}{response}})
}

