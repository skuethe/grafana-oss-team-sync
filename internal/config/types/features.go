package types

type Features struct {
	DisableFolders       bool `yaml:"disableFolders"`
	DisableUserSync      bool `yaml:"disableUserSync"`
	AddLocalAdminToTeams bool `yaml:"addLocalAdminToTeams"`
}

const (
	FeaturesDisableFoldersDefault  bool   = false
	FeaturesDisableFoldersFlagHelp string = "disable the folder sync feature"
	// FeaturesDisableFoldersFlagShort string = ""
	FeaturesDisableFoldersParameter string = "features.disableFolders"
	// FeaturesDisableFoldersVariable  string = "GOTS_"

	FeaturesDisableUsersDefault  bool   = false
	FeaturesDisableUsersFlagHelp string = "disable the user sync feature"
	// FeaturesDisableUsersFlagShort string = ""
	FeaturesDisableUsersParameter string = "features.disableUsers"
	// FeaturesDisableUsersVariable  string = "GOTS_"

	FeaturesAddLocalAdminToTeamsDefault  bool   = true
	FeaturesAddLocalAdminToTeamsFlagHelp string = "add the local Grafana admin user to each team you create"
	// FeaturesAddLocalAdminToTeamsFlagShort string = ""
	FeaturesAddLocalAdminToTeamsParameter string = "features.addLocalAdminToTeams"
	// FeaturesAddLocalAdminToTeamsVariable  string = "GOTS_"
)
