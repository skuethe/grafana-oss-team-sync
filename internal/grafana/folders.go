package grafana

import (
	"log/slog"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type folder struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	input  config.FolderSchema
	model  models.CreateFolderCommand
}

func (f *folder) doesFolderExist() bool {
	_, err := f.client.Folders.GetFolderByUID(f.model.UID)
	return err == nil
}

func (f *folder) createFolder() {
	_, err := f.client.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       f.model.Title,
		Description: f.model.Description,
		UID:         f.model.UID,
		ParentUID:   f.model.ParentUID,
	})
	if err != nil {
		f.log.Error("could not create Grafana folder", "error", err)
	} else {
		f.log.Info("created Grafana folder")
	}
}

func (f *folder) manageFolderPermissions() {

	var permissionList []*models.DashboardACLUpdateItem

	for teamName, teamPermission := range f.input.Permissions.Teams {
		team, err := f.client.Teams.SearchTeams(&teams.SearchTeamsParams{
			Name: &teamName,
		})
		if err != nil {
			f.log.Error("could not search for specified team - skipping", "error", err)
			continue
		}
		if len(team.Payload.Teams) == 0 {
			f.log.Error("skipping folder permissions for team", "team", teamName, "error", "team does not exist")
			continue
		}

		permerr := config.ValidateFolderPermission(teamPermission)
		if permerr != nil {
			f.log.Error("skipping folder permissions for team", "team", teamName, "error", permerr)
			continue
		}

		var permissionType models.PermissionType = models.PermissionType(teamPermission)

		permissionList = append(permissionList, &models.DashboardACLUpdateItem{
			Permission: permissionType,
			TeamID:     team.Payload.Teams[0].ID,
		})
	}

	_, err := f.client.FolderPermissions.UpdateFolderPermissions(f.model.UID, &models.UpdateDashboardACLCommand{
		Items: permissionList,
	})
	if err != nil {
		f.log.Error("could not update Grafana folder permissions", "error", err)
	} else {
		f.log.Info("Grafana folder permissions updated")
	}
}

func (g *GrafanaInstance) ProcessFolders() {

	foldersLog := slog.With(slog.String("package", "grafana.folders"))
	foldersLog.Info("processing Grafana folders")

	countSkipped := 0
	countCreated := 0

	folders := config.K.MapKeys("folders")

	for _, folderUID := range folders {
		var folderFromConfig config.FolderSchema

		config.K.Unmarshal("folders."+folderUID, &folderFromConfig)

		folderLog := slog.With(
			slog.Group("folder",
				slog.String("uid", folderUID),
				slog.String("title", folderFromConfig.Title),
			),
		)

		if len(folderFromConfig.Permissions.Teams) > 0 {
			config.K.MustInt64Map("folders." + folderUID + ".permissions.teams")
		}

		f := folder{
			client: g.api,
			log:    *folderLog,
			input:  folderFromConfig,
			model: models.CreateFolderCommand{
				UID:         strings.ToLower(folderUID),
				Title:       folderFromConfig.Title,
				Description: folderFromConfig.Description,
			},
		}

		if f.doesFolderExist() {
			countSkipped++
			folderLog.Debug("skipping Grafana folder because it already exists")
		} else {
			f.createFolder()
			countCreated++
		}

		f.manageFolderPermissions()
	}

	foldersLog.Info(
		"finished processing Grafana folders",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
