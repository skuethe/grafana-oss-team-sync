package grafana

import (
	"log/slog"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type UserType struct {
	client *goapi.GrafanaHTTPAPI
	log    slog.Logger
	form   models.AdminCreateUserForm
}

func (u *UserType) doesUserExist() bool {
	_, err := u.client.Users.GetUserByLoginOrEmail(u.form.Login)
	return err == nil
}

func (u *UserType) createUser() {
	_, err := u.client.AdminUsers.AdminCreateUser(&models.AdminCreateUserForm{
		Email:    u.form.Email,
		Login:    u.form.Login,
		Name:     u.form.Name,
		Password: u.form.Password,
	})
	if err != nil {
		u.log.Error("Could not create User", "error", err)
	} else {
		u.log.Info(
			"Created User",
			slog.Group("user",
				slog.String("login", u.form.Login),
				slog.String("email", u.form.Email),
			),
		)
	}
}

func ProcessUsers(c *goapi.GrafanaHTTPAPI, userList []models.AdminCreateUserForm) {
	usersLog := slog.With(slog.String("package", "grafana.users"))
	usersLog.Info("Processing Users")

	countSkipped := 0
	countCreated := 0

	for _, user := range userList {
		u := UserType{
			client: c,
			log:    *usersLog,
			form:   user,
		}
		if u.doesUserExist() {
			countSkipped++
			usersLog.Debug(
				"Skipped User",
				slog.Group("user",
					slog.String("login", user.Login),
					slog.String("email", user.Email),
				),
			)
		} else {
			u.createUser()
			countCreated++
		}
	}
	usersLog.Info(
		"Finished Users",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
