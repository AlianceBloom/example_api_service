package game

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"math"
)

const TournamentStatusInitial = "initial"
const TournamentStatusActive = "active"
const TournamentStatusFinished = "finished"

type Tournament struct {
	gorm.Model
	Uid    string `gorm:"type:varchar(100); unique_index; not null;"`
	Status        string `gorm:"varchar(16); default:'initial'"`
	MinimalDeposit int
	Players []*Player `gorm:"many2many:tournament_players;"`
	Bookings []*Booking
	game          *Game

	TournamentWinner *TournamentWinner
}

type TournamentWinner struct {
	gorm.Model

	Player *Player
	PlayerID uint

	Amount int
	Tournament *Tournament
	TournamentID uint
}

func (t *Tournament) JoinTournament(playerName string, LendersList... string) error {
	var lenders []*Player
	player := FindPlayer(playerName)
	if player == nil {
		return errors.New("Cant find player")
	}

	if Database.Where("player_id = ? AND tournament_id = ?", player.ID, t.ID).Table("tournament_players").RecordNotFound() {
		return errors.New("User already registred under tournament")
	}

	if len(LendersList) > 0 {
		Database.Where("uid IN (?)", LendersList).Find(&lenders)
	}

 	_, err := CreateBooking(t, player, lenders...)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tournament) FinishTournament(winnerName string) (*TournamentWinner, error) {
	if t.Status == TournamentStatusFinished {
		tournamentWinner := TournamentWinner{}
		if notFound := Database.Where("tournament_id = ?").Find(&tournamentWinner).RecordNotFound(); notFound {
			return nil, errors.New("TournamentWinner not found")
		}

		return &tournamentWinner, nil
	}

	player := FindPlayer(winnerName)
	if player == nil {
		return nil, errors.New("Player " + winnerName + " not found!" )
	}

	if Database.Where("tournament_id = ? AND player_id = ?", t.ID, player.ID).Find(&Booking{}).RecordNotFound() {
		return nil, errors.New("Player " + winnerName + " not for this tournament found!")
	}

	bookings := []*Booking{}
	lenders := []*Player{}
	winAmount := 0

	Database.Model(t).Preload("Lenders").Related(&bookings)

	tx := Database.Begin()

	for _, booking := range bookings {
		currentWin := booking.Amount

		if len(booking.Lenders) > 0 {
			currentWin *= len(booking.Lenders) + 1
		}
		if booking.PlayerID == player.ID {
			lenders = booking.Lenders
		}
		winAmount += currentWin
	}

	if len(lenders) > 0 {
		base := float64(winAmount) / float64(len(lenders) + 1)
		winAmount = int(math.Ceil(base))
	}

	var allPlayers     []*Player = append([]*Player{player}, lenders...)

	for _, currentPlayer := range allPlayers {

		err := tx.Exec("UPDATE players SET points = points + ? WHERE id = ?", winAmount, currentPlayer.ID ).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tournament := TournamentWinner{TournamentID: t.ID, Amount: winAmount, PlayerID: player.ID}
	err := tx.Create(&tournament).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Exec("UPDATE tournaments SET status = ? WHERE id = ?", TournamentStatusFinished, t.ID ).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &tournament, nil
}
