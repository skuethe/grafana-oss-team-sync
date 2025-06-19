package grafana

import (
	"log/slog"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/config/configtypes"
)

type Folder models.CreateFolderCommand

func (f *Folder) doesFolderExist() bool {
	_, err := Instance.api.Folders.GetFolderByUID(f.UID)
	return err == nil
}

func (f *Folder) createFolder() error {
	_, err := Instance.api.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       f.Title,
		Description: f.Description,
		UID:         f.UID,
		ParentUID:   f.ParentUID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *Folder) manageFolderPermissions(permissions configtypes.FolderPermissions) error {

	var permissionList []*models.DashboardACLUpdateItem

	for teamName, teamPermission := range permissions.Teams {
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

		// Validate defined permission for team
		if err := configtypes.ValidateGrafanaPermission(teamPermission); err != nil {
			slog.Error("skipping folder permissions for team",
				slog.String("team", teamName),
				slog.Any("error", err),
			)
			continue
		}

		permissionType := models.PermissionType(teamPermission)

		permissionList = append(permissionList, &models.DashboardACLUpdateItem{
			Permission: permissionType,
			TeamID:     team.Payload.Teams[0].ID,
		})
	}

	_, err := Instance.api.FolderPermissions.UpdateFolderPermissions(f.UID, &models.UpdateDashboardACLCommand{
		Items: permissionList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GrafanaInstance) ProcessFolders() {
	foldersLog := slog.With(slog.String("package", "grafana.folders"))

	if config.Instance.Features.DisableFolders {
		foldersLog.Info("folder feature disabled, skipping")
	} else if len(config.Instance.Folders) == 0 {
		foldersLog.Info("your folders input is empty, skipping")
	} else {
		foldersLog.Info("processing Grafana folders")

		countSkipped := 0
		countCreated := 0

		for folderUID, folder := range config.Instance.Folders {

			folderLog := slog.With(
				slog.Group("folder",
					slog.String("uid", folderUID),
					slog.String("title", folder.Title),
				),
			)

			f := Folder{
				UID:         strings.ToLower(folderUID),
				Title:       folder.Title,
				Description: folder.Description,
			}

			if f.doesFolderExist() {
				countSkipped++
				folderLog.Debug("skipping already existing Grafana folder")
			} else {
				if err := f.createFolder(); err != nil {
					folderLog.Error("could not create Grafana folder",
						slog.Any("error", err),
					)
				} else {
					folderLog.Info("created Grafana folder")
					countCreated++
				}
			}

			if err := f.manageFolderPermissions(folder.Permissions); err != nil {
				folderLog.Error("could not update Grafana folder permissions",
					slog.Any("error", err),
				)
			} else {
				folderLog.Info("Grafana folder permissions updated")
			}
		}

		foldersLog.Info("finished processing Grafana folders",
			slog.Group("folders",
				slog.Int("created", countCreated),
				slog.Int("existing", countSkipped),
			),
		)
	}
}
