package game

import (
	"testing"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestPlayer_TakePoints(t *testing.T) {

	validPlayer, err := testGame.AddPlayer(300, RandNumStr())

	_, err = validPlayer.TakePoints(40)

	if validPlayer.Points != 260 {
		t.Error("Player pounts must be decreesed")
	}

	if err != nil {
		t.Error("Player points decreasing, must be empty, got: ", err)
	}

	_, err = validPlayer.TakePoints(300)
	if err == nil {
		t.Error("Must return error, got points:", validPlayer.Points)
	}
}

func TestPlayer_FundPoints(t *testing.T) {
	validPlayer, err := testGame.AddPlayer(300, RandNumStr())
	validPlayer.FundPoints(33)
	if validPlayer.Points != 333 {
		t.Error("Player pounts must be decreesed")
	}
	if err != nil {
		t.Error("Player points decreasing, must be empty, got: ", err)
	}
}
