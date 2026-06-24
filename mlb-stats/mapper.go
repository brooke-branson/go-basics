package mlbstats

import (
	"fmt"
	"go-basics/models"
)

// statsResponse mirrors only the JSON fields we need from the MLB API.
type statsResponse struct {
	Stats []struct {
		Splits []struct {
			Season string `json:"season"`
			Stat   struct {
				AtBats     int `json:"atBats"`
				BaseOnBalls int `json:"baseOnBalls"`
				StrikeOuts int `json:"strikeOuts"`
				HomeRuns   int `json:"homeRuns"`
				Hits       int `json:"hits"`
				Doubles    int `json:"doubles"`
				Triples    int `json:"triples"`
			} `json:"stat"`
			Player struct {
				FullName string `json:"fullName"`
			} `json:"player"`
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
		} `json:"splits"`
	} `json:"stats"`
}

func mapStatsToPlayer(playerID string, payload statsResponse) (models.Player, error) {
	if len(payload.Stats) == 0 || len(payload.Stats[0].Splits) == 0 {
		return models.Player{}, fmt.Errorf("no hitting stats found for player %s", playerID)
	}

	split := payload.Stats[0].Splits[0]
	stat := split.Stat

	singles := stat.Hits - stat.Doubles - stat.Triples - stat.HomeRuns
	if singles < 0 {
		return models.Player{}, fmt.Errorf("invalid hit breakdown for player %s", playerID)
	}

	return models.Player{
		PlayerId:   playerID,
		PlayerName: split.Player.FullName,
		PlayerTeam: split.Team.Name,
		AtBats:     stat.AtBats,
		Walks:      stat.BaseOnBalls,
		Strikeouts: stat.StrikeOuts,
		HomeRuns:   stat.HomeRuns,
		Single:     singles,
		Double:     stat.Doubles,
		Triple:     stat.Triples,
	}, nil
}
