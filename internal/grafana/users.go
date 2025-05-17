package grafana

import (
	"crypto/rand"
	"log/slog"
	"math/big"

	"github.com/grafana/grafana-openapi-client-go/models"
)

type User models.AdminCreateUserForm
type Users []User

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

func (u *User) doesUserExist() bool {
	_, err := Instance.api.Users.GetUserByLoginOrEmail(u.Login)
	return err == nil
}

func (u *User) createUser() error {
	_, err := Instance.api.AdminUsers.AdminCreateUser(&models.AdminCreateUserForm{
		Email:    u.Email,
		Login:    u.Login,
		Name:     u.Name,
		Password: models.Password(generateSecurePassword()),
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) ProcessUsers() {
	usersLog := slog.With(slog.String("package", "grafana.users"))
	usersLog.Info("processing Grafana users")

	countSkipped := 0
	countCreated := 0

	for _, user := range *u {

		userLog := slog.With(
			slog.Group("user",
				slog.String("login", user.Login),
				slog.String("email", user.Email),
			),
		)
		if user.doesUserExist() {
			countSkipped++
			userLog.Debug("skipped Grafana user")
		} else {
			err := user.createUser()
			if err != nil {
				userLog.Error("could not create Grafana user", "error", err)
			} else {
				userLog.Info("created Grafana user")
				countCreated++
			}
		}
	}
	usersLog.Info(
		"finished processing Grafana users",
		slog.Group("stats",
			slog.Int("created", countCreated),
			slog.Int("skipped", countSkipped),
		),
	)
}
