package grafana

import (
	"log/slog"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type Folder struct {
	input config.FolderSchema
	model models.CreateFolderCommand
}

func (f *Folder) doesFolderExist() bool {
	_, err := Instance.api.Folders.GetFolderByUID(f.model.UID)
	return err == nil
}

func (f *Folder) createFolder() error {
	_, err := Instance.api.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       f.model.Title,
		Description: f.model.Description,
		UID:         f.model.UID,
		ParentUID:   f.model.ParentUID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *Folder) manageFolderPermissions() error {

	var permissionList []*models.DashboardACLUpdateItem

	for teamName, teamPermission := range f.input.Permissions.Teams {
		team, err := Instance.api.Teams.SearchTeams(&teams.SearchTeamsParams{
			Name: &teamName,
		})
		if err != nil {
			slog.Error("could not search for specified team - skipping",
				slog.Any("error", err),
			)
			continue
		}
		if len(team.Payload.Teams) == 0 {
			slog.Error("skipping folder permissions for team",
				slog.String("team", teamName),
				slog.String("error", "team does not exist"),
			)
			continue
		}

		permerr := config.ValidateGrafanaPermission(teamPermission)
		if permerr != nil {
			slog.Error("skipping folder permissions for team",
				slog.String("team", teamName),
				slog.Any("error", permerr),
			)
			continue
		}

		var permissionType models.PermissionType = models.PermissionType(teamPermission)

		permissionList = append(permissionList, &models.DashboardACLUpdateItem{
			Permission: permissionType,
			TeamID:     team.Payload.Teams[0].ID,
		})
	}

	_, err := Instance.api.FolderPermissions.UpdateFolderPermissions(f.model.UID, &models.UpdateDashboardACLCommand{
		Items: permissionList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GrafanaInstance) ProcessFolders() {

	foldersLog := slog.With(slog.String("package", "grafana.folders"))
	foldersLog.Info("processing Grafana folders")

	countSkipped := 0
	countCreated := 0

	folders := config.K.MapKeys(config.ConfigParamFolders)

	for _, folderUID := range folders {
		var folderFromConfig config.FolderSchema

		config.K.Unmarshal(config.ConfigParamFolders+"."+folderUID, &folderFromConfig)

		folderLog := slog.With(
			slog.Group("folder",
				slog.String("uid", folderUID),
				slog.String("title", folderFromConfig.Title),
			),
		)

		if len(folderFromConfig.Permissions.Teams) > 0 {
			config.K.MustInt64Map(config.ConfigParamFolders + "." + folderUID + ".permissions.teams")
		}

		f := Folder{
			input: folderFromConfig,
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
			err := f.createFolder()
			if err != nil {
				folderLog.Error("could not create Grafana folder",
					slog.Any("error", err),
				)
			} else {
				folderLog.Info("created Grafana folder")
				countCreated++
			}
		}

		err := f.manageFolderPermissions()
		if err != nil {
			folderLog.Error("could not update Grafana folder permissions",
				slog.Any("error", err),
			)
		} else {
			folderLog.Info("Grafana folder permissions updated")
		}
	}

	foldersLog.Info("finished processing Grafana folders",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
