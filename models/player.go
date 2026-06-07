package models

type Player struct {
	PlayerId   string `json:"player_id" bson:"player_id"`
	PlayerName string `json:"player_name" bson:"player_name"`
	PlayerTeam string `json:"player_team" bson:"player_team"`
	AtBats     int    `json:"at_bats" bson:"at_bats"`
	Walks      int    `json:"walks" bson:"walks"`
	Strikeouts int    `json:"strikeouts" bson:"strikeouts"`
	HomeRuns   int    `json:"home_runs" bson:"home_runs"`
	Single     int    `json:"single" bson:"single"`
	Double     int    `json:"double" bson:"double"`
	Triple     int    `json:"triple" bson:"triple"`
}