package game

import (
	"testing"
	"log"
)

func TestCreateBooking(t *testing.T) {
	tournament, err := testGame.NewTournament(RandNumStr(),190)

	if err != nil {
		log.Panic(err)
	}
	validPlayer, err := testGame.AddPlayer(300, RandNumStr())

	booking, err := CreateBooking(tournament, validPlayer)

	if err != nil {
		t.Error("Booking must be successfully generated, got:", err)
	}

	if booking.Amount != 190 {
		t.Error("Booking amount must be", tournament.MinimalDeposit, "got:", booking.Amount)
	}

	if len(booking.Lenders) != 0 {
		t.Error("Booking must be without lenders, got:", booking.Lenders)
	}

	if booking.Player.ID != validPlayer.ID {
		t.Error("Booking must be created under player, got: ", booking.Player)
	}

	_, err = CreateBooking(tournament, validPlayer)
	if err == nil {
		t.Error("Booking must not be generated, by dublication reason, got:", err)
	}
}

func TestCreateBooking2(t *testing.T) {

	tournament, err := testGame.NewTournament(RandNumStr(),200)

	if err != nil {
		log.Panic(err)
	}
	player, err := testGame.AddPlayer(150, RandNumStr())
	fundPlayer4, err := testGame.AddPlayer(300, RandNumStr())
	fundPlayer5, err := testGame.AddPlayer(300, RandNumStr())
	fundPlayer6, err := testGame.AddPlayer(300, RandNumStr())

	booking, err := CreateBooking(tournament, player, fundPlayer4, fundPlayer5, fundPlayer6)
	if err != nil {
		t.Error("Booking must be successfull generated, got:", err)
	}

	if booking.Amount != 50 {
		t.Error("Booking amount must be", tournament.MinimalDeposit, "got:", booking.Amount)
	}

	if count := testDB.Model(booking).Association("Lenders").Count(); count != 3 {
		t.Error("Booking must be with", 3,"lenders, got:", count)
	}

	if booking.Player.ID != player.ID {
		t.Error("Booking must be created under player, got: ", booking.Player)
	}

	_, err = CreateBooking(tournament, player)
	if err == nil {
		t.Error("Booking must not be generated, by dublication reason, got", err)
	}

	if fundPlayer4.Points != 250 {
		t.Error("Each player must have", 250, "points, got:", fundPlayer4.Points )
	}
	if fundPlayer5.Points != 250 {
		t.Error("Each player must have", 250, "points, got:", fundPlayer5.Points )
	}
	if fundPlayer6.Points != 250 {
		t.Error("Each player must have", 250, "points, got:", fundPlayer6.Points )
	}
}

