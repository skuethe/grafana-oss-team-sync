package grafana

import (
	"log"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func doesUserExist(c *goapi.GrafanaHTTPAPI, login string) bool {
	_, err := c.Users.GetUserByLoginOrEmail(login)
	return err == nil
}

func createUser(c *goapi.GrafanaHTTPAPI, user models.AdminCreateUserForm) {
	_, err := c.AdminUsers.AdminCreateUser(&models.AdminCreateUserForm{
		Email:    user.Email,
		Login:    user.Login,
		Name:     user.Name,
		Password: user.Password,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("  Created: \"%v\" (%v)", user.Login, user.Email)
	}
}

func ProcessUsers(c *goapi.GrafanaHTTPAPI, userList []models.AdminCreateUserForm) {
	log.SetPrefix("[Grafana.Users] ")
	log.Println("Processing ...")

	countSkipped := 0
	countCreated := 0

	for _, user := range userList {
		if doesUserExist(c, user.Login) {
			countSkipped++
		} else {
			createUser(c, user)
			countCreated++
		}
	}
	log.Printf("Created: %v; Skipped: %v", countCreated, countSkipped)
}
