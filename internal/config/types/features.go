package types

type Features struct {
	DisableFolders       bool `yaml:"disableFolders"`
	DisableUserSync      bool `yaml:"disableUserSync"`
	AddLocalAdminToTeams bool `yaml:"addLocalAdminToTeams"`
}

const (
	// FeaturesDisableFoldersFlagShort string = ""
	// FeaturesDisableFoldersVariable  string = "GOTS_"
	FeaturesDisableFoldersDefault   bool   = false
	FeaturesDisableFoldersFlagHelp  string = "disable the folder sync feature"
	FeaturesDisableFoldersParameter string = "features.disableFolders"

	// FeaturesDisableUsersFlagShort string = ""
	// FeaturesDisableUsersVariable  string = "GOTS_"
	FeaturesDisableUsersDefault   bool   = false
	FeaturesDisableUsersFlagHelp  string = "disable the user sync feature"
	FeaturesDisableUsersParameter string = "features.disableUsers"

	// FeaturesAddLocalAdminToTeamsFlagShort string = ""
	// FeaturesAddLocalAdminToTeamsVariable  string = "GOTS_"
	FeaturesAddLocalAdminToTeamsDefault   bool   = true
	FeaturesAddLocalAdminToTeamsFlagHelp  string = "add the local Grafana admin user to each team you create"
	FeaturesAddLocalAdminToTeamsParameter string = "features.addLocalAdminToTeams"
)
