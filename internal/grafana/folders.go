package grafana

import (
	"log/slog"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type FolderType struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	form   models.CreateFolderCommand
}

func (f *FolderType) doesFolderExist() bool {
	_, err := f.client.Folders.GetFolderByUID(f.form.UID)
	return err == nil
}

func (f *FolderType) createFolder() {
	_, err := f.client.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       f.form.Title,
		Description: f.form.Description,
		UID:         f.form.UID,
		ParentUID:   f.form.ParentUID,
	})
	if err != nil {
		f.log.Error(err.Error())
	} else {
		f.log.Info(
			"Created Grafana Folder",
			slog.Group("folder",
				slog.String("uid", f.form.UID),
				slog.String("title", f.form.Title),
				slog.String("description", f.form.Description),
			),
		)
	}
}

func (g *GrafanaInstance) ProcessFolders(folderList *[]models.CreateFolderCommand) {
	foldersLog := slog.With(slog.String("package", "grafana.folders"))
	foldersLog.Info("Processing Grafana Folders")

	countSkipped := 0
	countCreated := 0

	for _, folder := range *folderList {
		f := FolderType{
			client: g.api,
			log:    *foldersLog,
			form:   folder,
		}
		if f.doesFolderExist() {
			countSkipped++
			foldersLog.Debug(
				"Skipped Grafana Folder",
				slog.Group("user",
					slog.String("uid", folder.UID),
					slog.String("title", folder.Title),
					slog.String("description", folder.Description),
				),
			)
		} else {
			f.createFolder()
			countCreated++
		}
	}
	foldersLog.Info(
		"Finished Grafana Folders",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
