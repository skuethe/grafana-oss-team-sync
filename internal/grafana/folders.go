package grafana

import (
	"log/slog"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type configFolderSchema struct {
	Title       string `koanf:"title"`
	Description string `koanf:"description"`
	Permissions struct {
		Teams map[string]permissionsAllowed `koanf:"teams"`
	} `koanf:"permissions"`
}

type folder struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	input  configFolderSchema
	model  models.CreateFolderCommand
}

type permissionsAllowed int64

const (
	PermissionViewer permissionsAllowed = 1 << iota
	PermissionEditor
	PermissionAdmin
)

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
		f.log.Error("Could not create Grafana Folder", "error", err)
	} else {
		f.log.Info("Created Grafana Folder")
	}
}

func (f *folder) manageFolderPermissions() {

	var permissionList []*models.DashboardACLUpdateItem

	for teamName, teamPermission := range f.input.Permissions.Teams {
		team, err := f.client.Teams.SearchTeams(&teams.SearchTeamsParams{
			Name: &teamName,
		})
		if err != nil {
			f.log.Error("Could not search for specific team - skipping", "error", err)
			continue
		}
		if len(team.Payload.Teams) == 0 {
			f.log.Warn("Specified Team for folder permissions does not exist - skipping")
			continue
		}

		switch teamPermission {
		case PermissionViewer, PermissionEditor, PermissionAdmin:
		default:
			f.log.Warn("Permission not allowed - skipping", "wrong.permission", teamPermission)
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
		f.log.Error("Could not update Grafana Folder Permissions", "error", err)
	} else {
		f.log.Info("Grafana Folder Permissions updated")
	}
}

func (g *GrafanaInstance) processFolders() {

	foldersLog := slog.With(slog.String("package", "grafana.folders"))
	foldersLog.Info("Initializing Grafana Folders")

	countSkipped := 0
	countCreated := 0

	folders := config.K.MapKeys("folders")

	for _, folderUID := range folders {
		var folderFromConfig configFolderSchema

		config.K.Unmarshal("folders."+folderUID, &folderFromConfig)

		folderLog := slog.With(
			slog.Group("folder",
				slog.String("uid", folderUID),
				slog.String("title", folderFromConfig.Title),
				slog.String("description", folderFromConfig.Description),
				slog.Any("permissions", config.K.Get("folders."+folderUID+".permissions")),
			),
		)
		folderLog.Info("Processing Grafana Folder")

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
			folderLog.Debug("Grafana Folder already exists - skipping")
		} else {
			f.createFolder()
			countCreated++
		}

		f.manageFolderPermissions()
	}

	foldersLog.Info(
		"Finished Grafana Folders",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
