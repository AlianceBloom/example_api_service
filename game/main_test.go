package game

import (
	"testing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"io/ioutil"
	"os"
	"time"
	"math/rand"
	"strconv"
)

var testDB *gorm.DB
var testGame *Game

func RandNumStr() string {
	rand.Seed(time.Now().UTC().UnixNano())
	randNumb := rand.Intn(5000)
	return strconv.Itoa(randNumb)
}

func TestMain(m *testing.M) {
	tmpDir, _ := ioutil.TempDir("", "text_game_engine")
	tmpDbPath := tmpDir + "/" + RandNumStr() + "testDb.sqlite"
	log.Println("DBpath: ", tmpDbPath)

	var err error
	testDB, err = gorm.Open("sqlite3", tmpDbPath)

	if err != nil {
		log.Panic("DB connection error:", err)
	}

	if err != nil {
		log.Panic("Drop DB err:", err)
	}
	defer testDB.Close()
	defer os.RemoveAll(tmpDir)
	
	testGame = InitializeGame(testDB)
	os.Exit(m.Run())
}


func TestGame_AddPlayer(t *testing.T) {
	validPlayer, err := testGame.AddPlayer(300, RandNumStr())

	if err != nil {
		t.Error("Player must be returned, got: ", validPlayer)
	}

	notValidPlayer, err := testGame.AddPlayer(300, validPlayer.Uid)
	if err == nil {
		t.Error("Must return an error about player dublication", notValidPlayer.Uid)
	}
}

func TestInitializeGame(t *testing.T)  {
	tournament, err := testGame.NewTournament(RandNumStr(), 300)
	if err != nil {
		t.Error("Must create a new tournament, got:", err)
	}
	if tournament.ID == 0 {
		t.Error("Must create a testDB record for tournament, got", tournament.ID)
	}
	if tournament.MinimalDeposit != 300 {
		t.Error("Must stote a minimal deposit, got", tournament.MinimalDeposit)
	}
	if tournament.Status != TournamentStatusActive {
		t.Error("Must be created with active status, got", tournament.Status)
	}
}

