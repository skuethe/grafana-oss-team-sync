// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

type Features struct {
	DisableFolders       bool `yaml:"disableFolders"`
	DisableUserSync      bool `yaml:"disableUserSync"`
	AddLocalAdminToTeams bool `yaml:"addLocalAdminToTeams"`
}

const (
	// FeaturesAddLocalAdminToTeamsFlagShort string = ""
	FeaturesAddLocalAdminToTeamsDefault   bool   = true
	FeaturesAddLocalAdminToTeamsFlagHelp  string = "feature: add the local Grafana admin user to each team you create"
	FeaturesAddLocalAdminToTeamsParameter string = "features.addLocalAdminToTeams"
	FeaturesAddLocalAdminToTeamsOptimized string = "addlocaladmintoteams"

	// FeaturesDisableFoldersFlagShort string = ""
	FeaturesDisableFoldersDefault   bool   = false
	FeaturesDisableFoldersFlagHelp  string = "feature: disable folders"
	FeaturesDisableFoldersParameter string = "features.disableFolders"
	FeaturesDisableFoldersOptimized string = "disablefolders"

	// FeaturesDisableUsersFlagShort string = ""
	FeaturesDisableUsersDefault   bool   = false
	FeaturesDisableUsersFlagHelp  string = "feature: disable the user sync"
	FeaturesDisableUsersParameter string = "features.disableUserSync"
	FeaturesDisableUsersOptimized string = "disableusersync"
)
