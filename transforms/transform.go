package transforms

import (
	"go-basics/models"
"encoding/xml"
)

type TransformedPlayer struct {
	XMLName           xml.Name  `xml:"Player" json:"-"`
	PlayerId          string  `xml:"player_id" json:"player_id"`
	PlayerName        string  `xml:"player_name" json:"player_name"`
	PlayerTeam        string  `xml:"player_team" json:"player_team"`
	BattingAverage    float64 `xml:"batting_average" json:"batting_average"`
	HomeRunPercentage float64 `xml:"home_run_percentage" json:"home_run_percentage"`
	TotalBases        int     `xml:"total_bases" json:"total_bases"`
	SluggingPercentage float64 `xml:"slugging_percentage"`
	OnBasePlusSlugging float64 `xml:"on_base_plus_slugging" json:"on_base_plus_slugging"`
	OnBasePercentage    float64 `xml:"on_base_percentage" json:"on_base_percentage"`
}

func TransformPlayer(player models.Player) TransformedPlayer {
	hits := player.Single + player.Double + player.Triple + player.HomeRuns
	battingAverage := float64(hits) / float64(player.AtBats)
	homeRunPercentage := float64(player.HomeRuns) / float64(player.AtBats)
	totalBases := player.Single + 2 * player.Double + 3 * player.Triple + 4 * player.HomeRuns
	sluggingPercentage := float64(totalBases) / float64(player.AtBats)
	onBasePercentage := float64(hits) / float64(player.AtBats + player.Walks) // TODO: Add correct formula
	onBasePlusSlugging := onBasePercentage + sluggingPercentage

	return TransformedPlayer{
		XMLName:           xml.Name{Space: "", Local: "Player"},
		PlayerId:          player.PlayerId,
		PlayerName:        player.PlayerName,
		PlayerTeam:        player.PlayerTeam,
		BattingAverage:    battingAverage,
		HomeRunPercentage: homeRunPercentage,
		TotalBases:        totalBases,
		SluggingPercentage: sluggingPercentage,
		OnBasePlusSlugging: onBasePlusSlugging,
		OnBasePercentage: onBasePercentage,
	}
}
