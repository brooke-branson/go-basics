package models

type Player struct {
	PlayerId string `json:"player_id"`
	PlayerName string `json:"player_name"`
	PlayerTeam string `json:"player_team"`
}