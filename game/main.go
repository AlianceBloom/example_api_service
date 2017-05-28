package game

import (
	"github.com/jinzhu/gorm"
	"errors"
)

var Database *gorm.DB

type Game struct {
	database *gorm.DB
}


func InitializeGame(DB *gorm.DB) *Game {
	Database = DB
	game := Game{}
	Database.AutoMigrate(Player{}, Tournament{}, Booking{}, TournamentWinner{})

	Database.FirstOrCreate(&Player{}, &Player{ Uid: "P1", Points: 0})
	Database.FirstOrCreate(&Player{}, &Player{ Uid: "P2", Points: 0})
	Database.FirstOrCreate(&Player{}, &Player{ Uid: "P3", Points: 0})
	Database.FirstOrCreate(&Player{}, &Player{ Uid: "P4", Points: 0})
	Database.FirstOrCreate(&Player{}, &Player{ Uid: "P5", Points: 0})

	return &game
}

func FindPlayer(playerName string) *Player {
	player := Player{}

	if notFund := Database.Where("uid = ?", playerName).Find(&player).RecordNotFound(); notFund {
		return nil
	}

	return &player
}

func FindTournament(tournamentID interface{}) *Tournament {
	tournament := Tournament{}
	if notFound := Database.Where("uid = ?", tournamentID).Find(&tournament).RecordNotFound(); notFound {
		return nil
	}
	return &tournament
}


func (g *Game) NewTournament(Uid string, deposit int) (*Tournament, error) {
	if tournament := FindTournament(Uid); tournament != nil {
		return nil, errors.New("Tournament already exist")
	}

	tournament := Tournament{MinimalDeposit: deposit, Status: TournamentStatusActive, Uid: Uid}

	err := Database.Create(&tournament).GetErrors()
	if len(err) == 0 {
		return &tournament, nil
	}

	return nil, errors.New("Error while tournament creation")
}

func (g *Game) AddPlayer(initialPoints int, Uid string) (*Player, error) {
	if player := FindPlayer(Uid); player != nil {
		return nil, errors.New("Player already exist")
	}

	player := Player{Points: initialPoints, Uid: Uid}

	err := Database.Create(&player).GetErrors()
	if len(err) == 0 {
		return &player, nil
	}


	return nil, errors.New("Error while player create")
}

