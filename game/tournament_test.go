package game

import (
	"testing"
)

func TestTournament_JoinTournament(t *testing.T) {
	player, _:= testGame.AddPlayer(300, RandNumStr())
	tournament, err := testGame.NewTournament(RandNumStr(),150)


	err = tournament.JoinTournament(player.Uid)
	if err != nil {
		t.Error("Must add player to tournament, got", player.Uid)
	}

	if len(tournament.Players) == 0 || tournament.Players[0].Uid != player.Uid {
		t.Error("Player must be assigned to tournament")
	}
}

func TestTournament_JoinTournament2(t *testing.T) {
	player, _ := testGame.AddPlayer(300, RandNumStr())
	player2, _ := testGame.AddPlayer(300, RandNumStr())
	player3, _ := testGame.AddPlayer(300, RandNumStr())
	player4, _ := testGame.AddPlayer(300, RandNumStr())
	tournament, err := testGame.NewTournament(RandNumStr(),150)


	err = tournament.JoinTournament(player.Uid, player2.Uid, player3.Uid, player4.Uid)
	if err != nil {
		t.Error("Must add player to tournament, got", player.Uid)
	}
	if len(tournament.Players) == 0 || tournament.Players[0].Uid != player.Uid {
		t.Error("Player must be assigned to tournament")
	}

	if count := testDB.Model(&tournament).Association("Bookings").Count(); count > 1 {
		t.Error("Must have 1 booking, got",  count)
	}
}

func TestTournament_JoinTournament3(t *testing.T) {
	player, _ := testGame.AddPlayer(130, RandNumStr())
	player2, _ := testGame.AddPlayer(300, RandNumStr())
	player3, _ := testGame.AddPlayer(300, RandNumStr())
	player4, _ := testGame.AddPlayer(300, RandNumStr())
	tournament, err := testGame.NewTournament(RandNumStr(), 150)


	err = tournament.JoinTournament(player.Uid, player2.Uid, player3.Uid, player4.Uid)
	if err != nil {
		t.Error("Must add player to tournament, got", err)
	}

	if len(tournament.Players) == 0 || tournament.Players[0].Uid != player.Uid {
		t.Error("Player must be assigned to tournament")
	}

	testDB.Model(tournament).Related(&tournament.Bookings)
	if len(tournament.Bookings) == 0 || len(tournament.Bookings) > 1 {
		t.Error("Must have 1 booking, got",  len(tournament.Bookings))
	}

	count := 0
	testDB.Table("player_booking_lenders").Where("booking_id = ?", tournament.Bookings[0].ID).Count(&count)
	if count != 3 {
		t.Error("Lenders slice must be equal to 3, got:", count)
	}
}

func TestGame_FinishTournament(t *testing.T) {
	player, _ := testGame.AddPlayer(160, RandNumStr())
	player2, _ := testGame.AddPlayer(169, RandNumStr())
	tournament, err := testGame.NewTournament(RandNumStr(),150)

	tournament.JoinTournament(player.Uid, player2.Uid)

	_, err = tournament.FinishTournament(player.Uid)
	if err != nil {
		t.Error("Must not return any error, got: ", err)
	}

	_, err = tournament.FinishTournament("TestName")
	if err == nil {
		t.Error("Must return error about missing player")
	}

	_, err = tournament.FinishTournament(player2.Uid)
	if err == nil {
		t.Error("Must return error about missing player")
	}
}

func TestTournament_FinishTournament2(t *testing.T) {
	player1, _ := testGame.AddPlayer(300, "Test-P1")
	player2, _ := testGame.AddPlayer(300, "Test-P2")
	player3, _ := testGame.AddPlayer(300, "Test-P3")
	player4, _ := testGame.AddPlayer(500, "Test-P4")
	player5, _ := testGame.AddPlayer(1000, "Test-P5")

	tournament, _ := testGame.NewTournament(RandNumStr(), 1000)
	tournament.JoinTournament(player1.Uid, player2.Uid, player3.Uid, player4.Uid)
	tournament.JoinTournament(player5.Uid)

	tournament.FinishTournament(player1.Uid)

	if testDB.Find(&player1); player1.Points != 550 {
		t.Error("Player 1 must contain 550, got", player1.Points)
	}

	if testDB.Find(&player2); player2.Points != 550 {
		t.Error("Player 2 must contain 550, got", player2.Points)
	}

	if testDB.Find(&player3); player3.Points != 550 {
		t.Error("Player 3 must contain 550, got", player3.Points)
	}

	if testDB.Find(&player4); player4.Points != 750 {
		t.Error("Player 4 must contain 750, got", player4.Points)
	}

	if testDB.Find(&player5); player5.Points != 0 {
		t.Error("Player 5 must contain 0, got", player5.Points)
	}

	if count := testDB.Model(tournament).Association("TournamentWinner").Count(); count != 1 {
		t.Error("Tournament winner must be created ", count)
	}
}
