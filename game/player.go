package game

import (
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
)

type Player struct {
	gorm.Model

	Uid    string `gorm:"type:varchar(100); unique_index; not null;"`
	Points int
	Tournaments []Tournament `gorm:"many2many:player_tournaments;"`
}


func (p *Player) TakePoints(amount int) (int, error) {
	if checkAmount := p.Points - amount; checkAmount < 0 {
		return 0, errors.New("Not enough points")
	}

	tx := Database.Begin()
	if err := tx.Exec("UPDATE players SET points = points - ? WHERE id = ?", amount, p.ID).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	p.Points -= amount
	tx.Commit()

	return p.Points, nil
}


func (p *Player) FundPoints(amount int) (int, error) {

	tx := Database.Begin()
	if err := tx.Exec("UPDATE players SET points = points + ? WHERE id = ?", amount, amount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	p.Points += amount
	tx.Commit()

	return p.Points, nil
}
