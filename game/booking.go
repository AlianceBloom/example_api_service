package game

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"math"

)

type Booking struct {
	gorm.Model

	Amount int

	PlayerID uint
	Player *Player

	TournamentID uint
	Tournament *Tournament

	Lenders []*Player `gorm:"many2many:player_booking_lenders;"`
}

func CreateBooking(tournament *Tournament, player *Player, lenders... *Player) (*Booking, error) {
	var allPlayers      []*Player = append([]*Player{player}, lenders...)
	var requestedAmount int       = calculateRequestedAmount(tournament, allPlayers...)
	var booking         Booking = Booking{Player: player, Tournament: tournament, Amount: requestedAmount}

	tx := Database.Begin()

	if player.Points < tournament.MinimalDeposit {
		booking.Lenders = lenders
		booking.Amount = requestedAmount
	}

	for _, currentPlayer := range allPlayers {
		err := tx.Model(currentPlayer).UpdateColumn("points", gorm.Expr("points - ?", requestedAmount)).Error
		if err != nil || currentPlayer.Points < 0 {
			tx.Rollback()
			return nil, errors.New("Can reserve requested point amount from player: " + currentPlayer.Uid)
		}
		currentPlayer.Points -= requestedAmount
	}


	err := tx.Create(&booking).Error;
	if err != nil {
		tx.Rollback()
		return nil, errors.New("Cant create booking")
	}

	tx.Commit()
	return &booking, nil
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	var count int = 0

	Database.Table("bookings").Where("player_id = ? AND tournament_id = ?", b.Player.ID, b.Tournament.ID).Count(&count)
	if count > 0 {
		return errors.New("Booking already exist")
	}

	return nil
}

func (b *Booking) AfterCreate(tx *gorm.DB) (err error) {
	err = tx.Model(b.Tournament).Association("Players").Append([]*Player{b.Player}).Error
	if err != nil {
		return err
	}

	return nil
}

func calculateRequestedAmount(t *Tournament, playersList... *Player) int {
	var division float64 = float64(t.MinimalDeposit) / float64(len(playersList))
	return int(math.Ceil(division))
}
