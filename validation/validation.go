package validation

import (
	"go-basics/models"
	"errors"
)

func ValidatePlayer(player models.Player) error {
	if player.PlayerId == "" {
		return errors.New("player id is required")
	}
	if player.PlayerName == "" {
		return errors.New("player name is required")
	}
	if player.PlayerTeam == "" {
		return errors.New("player team is required")
	}
	if player.AtBats <= 0 {
		return errors.New("at bats is required")
	}
	if player.Walks < 0 {
		return errors.New("walks cannot be negative")
	}
	if player.Strikeouts < 0 {
		return errors.New("strikeouts cannot be negative")
	}
	if player.HomeRuns < 0 {
		return errors.New("home runs cannot be negative")
	}
	if player.Single < 0 {
		return errors.New("single cannot be negative")
	}
	if player.Double < 0 {
		return errors.New("double cannot be negative")
	}
	if player.Triple < 0 {
		return errors.New("triple cannot be negative")
	}
	return nil
}