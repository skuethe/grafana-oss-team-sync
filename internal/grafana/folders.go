package grafana

import (
	"log"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func doesFolderExist(c *goapi.GrafanaHTTPAPI, uid string) bool {
	_, err := c.Folders.GetFolderByUID(uid)
	return err == nil
}

func createFolder(c *goapi.GrafanaHTTPAPI, folder models.CreateFolderCommand) {
	_, err := c.Folders.CreateFolder(&models.CreateFolderCommand{
		Title:       folder.Title,
		Description: folder.Description,
		UID:         folder.UID,
		ParentUID:   folder.ParentUID,
	})
	if err != nil {
		log.Println("Whoops")
		log.Fatal(err)
	} else {
		log.Printf("  Created: \"%v\"", folder.UID)
	}
}

func ProcessFolders(c *goapi.GrafanaHTTPAPI, folderList []models.CreateFolderCommand) {
	log.SetPrefix("[Grafana.Folders] ")
	log.Println("Processing ...")

	countSkipped := 0
	countCreated := 0

	for _, folder := range folderList {
		if doesFolderExist(c, folder.UID) {
			countSkipped++
		} else {
			createFolder(c, folder)
			countCreated++
		}
	}
	log.Printf("Created: %v; Skipped: %v", countCreated, countSkipped)
}
