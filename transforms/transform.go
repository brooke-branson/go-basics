package transforms

type RawPlayer struct {
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	PlayerTeam string `json:"player_team"`
	AtBats     int    `json:"at_bats"`
	Walks      int    `json:"walks"`
	Strikeouts int    `json:"strikeouts"`
	HomeRuns   int    `json:"home_runs"`
	Single     int    `json:"single"`
	Double     int    `json:"double"`
	Triple     int    `json:"triple"`
	HomeRun    int    `json:"home_run"`
}

type TransformedPlayer struct {
	XmlName           string  `xml:"Player" json:"-"`
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

func TransformPlayer(raw RawPlayer) TransformedPlayer {
	hits := raw.Single + raw.Double + raw.Triple + raw.HomeRun
	battingAverage := float64(hits) / float64(raw.AtBats)
	homeRunPercentage := float64(raw.HomeRuns) / float64(raw.AtBats)
	totalBases := raw.Single + 2 * raw.Double + 3 * raw.Triple + 4 * raw.HomeRun
	sluggingPercentage := float64(totalBases) / float64(raw.AtBats)
	onBasePercentage := float64(hits) / float64(raw.AtBats + raw.Walks) // TODO: Add correct formula
	onBasePlusSlugging := onBasePercentage + sluggingPercentage

	return TransformedPlayer{
		XmlName:           "Player",
		PlayerId:          raw.PlayerId,
		PlayerName:        raw.PlayerName,
		PlayerTeam:        raw.PlayerTeam,
		BattingAverage:    battingAverage,
		HomeRunPercentage: homeRunPercentage,
		TotalBases:        totalBases,
		SluggingPercentage: sluggingPercentage,
		OnBasePlusSlugging: onBasePlusSlugging,
		OnBasePercentage: onBasePercentage,
	}
}
