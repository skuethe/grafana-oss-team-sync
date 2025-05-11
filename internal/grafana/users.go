package grafana

import (
	"crypto/rand"
	"log/slog"
	"math/big"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type user struct {
	client *client.GrafanaHTTPAPI
	log    slog.Logger
	form   models.AdminCreateUserForm
}

func generateSecurePassword() string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	const length = 32

	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := range password {
		index, _ := rand.Int(rand.Reader, charsetLength)
		password[i] = charset[index.Int64()]
	}

	return string(password)
}

func (u *user) doesUserExist() bool {
	_, err := u.client.Users.GetUserByLoginOrEmail(u.form.Login)
	return err == nil
}

func (u *user) createUser() {
	_, err := u.client.AdminUsers.AdminCreateUser(&models.AdminCreateUserForm{
		Email:    u.form.Email,
		Login:    u.form.Login,
		Name:     u.form.Name,
		Password: models.Password(generateSecurePassword()),
	})
	if err != nil {
		u.log.Error("Could not create User", "error", err)
	} else {
		u.log.Info("Created Grafana User")
	}
}

func (g *GrafanaInstance) processUsers(userList *[]models.AdminCreateUserForm) {
	usersLog := slog.With(slog.String("package", "grafana.users"))
	usersLog.Info("Processing Grafana Users")

	countSkipped := 0
	countCreated := 0

	for _, instance := range *userList {

		userLog := slog.With(
			slog.Group("user",
				slog.String("login", instance.Login),
				slog.String("email", instance.Email),
			),
		)
		userLog.Info("Processing Grafana User")

		u := user{
			client: g.api,
			log:    *userLog,
			form:   instance,
		}
		if u.doesUserExist() {
			countSkipped++
			userLog.Debug(
				"Skipped Grafana User",
			)
		} else {
			u.createUser()
			countCreated++
		}
	}
	usersLog.Info(
		"Finished Grafana Users",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
