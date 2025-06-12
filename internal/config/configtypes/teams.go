package configtypes

type Teams []string

const (
	// TeamsFlagShort string = ""
	TeamsDefault   string = ""
	TeamsFlagHelp  string = "the comma-seperated list of teams you want to sync"
	TeamsParameter string = "teams"
	TeamsVariable  string = "GOTS_TEAMS"
)
