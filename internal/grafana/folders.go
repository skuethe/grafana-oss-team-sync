package grafana

import (
	"log/slog"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type ConfigFolderSchema struct {
	Title       string `koanf:"title"`
	Description string `koanf:"description"`
	Permissions struct {
		Teams map[string]permissionsAllowed `koanf:"teams"`
	} `koanf:"permissions"`
}

type Folder struct {
	Client *client.GrafanaHTTPAPI
	Log    slog.Logger
	Input  ConfigFolderSchema
	Model  models.CreateFolderCommand
}

type permissionsAllowed int64

const (
	PermissionViewer permissionsAllowed = 1 << iota
	PermissionEditor
	PermissionAdmin
)

func (f *Folder) doesFolderExist() bool {
	_, err := f.Client.Folders.GetFolderByUID(f.Model.UID)
	return err == nil
}

func (f *Folder) createFolder() {
	_, err := f.Client.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       f.Model.Title,
		Description: f.Model.Description,
		UID:         f.Model.UID,
		ParentUID:   f.Model.ParentUID,
	})
	if err != nil {
		f.Log.Error("Could not create Grafana Folder", "error", err)
	} else {
		f.Log.Info("Created Grafana Folder")
	}
}

func (f *Folder) manageFolderPermissions() {

	var permissionList []*models.DashboardACLUpdateItem

	for key, val := range f.Input.Permissions.Teams {
		team, err := f.Client.Teams.SearchTeams(&teams.SearchTeamsParams{
			Name: &key,
		})
		if err != nil {
			f.Log.Error("Could not search for specific team - skipping", "error", err)
			continue
		}
		if len(team.Payload.Teams) == 0 {
			f.Log.Warn("Specified Team for folder permissions does not exist - skipping")
			continue
		}

		switch val {
		case PermissionViewer, PermissionEditor, PermissionAdmin:
		default:
			f.Log.Warn("Permission not allowed - skipping", "wrong.permission", val)
			continue
		}

		var permissionType models.PermissionType = models.PermissionType(val)

		// validationError := permissionType.Validate(strfmt.Default)
		// if validationError != nil {
		// 	slog.Error("Could not validate permission input", "input", val, "error", validationError)
		// 	break
		// }

		permissionList = append(permissionList, &models.DashboardACLUpdateItem{
			Permission: permissionType,
			TeamID:     team.Payload.Teams[0].ID,
		})
	}

	_, err := f.Client.FolderPermissions.UpdateFolderPermissions(f.Model.UID, &models.UpdateDashboardACLCommand{
		Items: permissionList,
	})
	if err != nil {
		f.Log.Error("Could not update Grafana Folder Permissions", "error", err)
	} else {
		f.Log.Info("Grafana Folder Permissions updated")
	}
}

func (g *GrafanaInstance) processFolders() {

	foldersLog := slog.With(slog.String("package", "grafana.folders"))
	foldersLog.Info("Initializing Grafana Folders")

	countSkipped := 0
	countCreated := 0

	folders := config.K.MapKeys("folders")

	for _, folderUID := range folders {
		var folderFromConfig ConfigFolderSchema

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

		f := Folder{
			Client: g.api,
			Log:    *folderLog,
			Input:  folderFromConfig,
			Model: models.CreateFolderCommand{
				UID:         strings.ToLower(folderUID),
				Title:       folderFromConfig.Title,
				Description: folderFromConfig.Description,
			},
		}

		if f.doesFolderExist() {
			countSkipped++
			foldersLog.Debug("Grafana Folder already exists - skipping")
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
