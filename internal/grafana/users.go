package grafana

import (
	"crypto/rand"
	"log/slog"
	"math/big"
	"slices"

	"github.com/grafana/grafana-openapi-client-go/client/users"
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

func (u *User) searchUser() (*users.GetUserByLoginOrEmailOK, error) {
	result, err := Instance.api.Users.GetUserByLoginOrEmail(u.Login)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *User) doesUserExist() bool {
	_, err := u.searchUser()
	return err == nil
}

func (u *User) getUserID() (*int64, error) {
	result, err := u.searchUser()
	if err != nil {
		return nil, err
	}
	return &result.Payload.ID, nil
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

func (t *Teams) ProcessUsers() {
	usersLog := slog.With(slog.String("package", "grafana.users"))
	usersLog.Info("processing Grafana users")

	countCreated := 0
	countDuplicate := 0
	countSkipped := 0

	var globalUserList *Users = &Users{}

	for _, team := range *t {
		for _, user := range *team.Users {
			userExists := false
			if slices.Contains(*globalUserList, user) {
				userExists = true
				usersLog.Debug("skipping duplicate user")
				countDuplicate++
			}
			if !userExists {
				*globalUserList = append(*globalUserList, user)
			}
		}
	}

	for _, user := range *globalUserList {

		userLog := slog.With(
			slog.Group("user",
				slog.String("login", user.Login),
				slog.String("email", user.Email),
			),
		)
		if user.doesUserExist() {
			countSkipped++
			userLog.Debug("skipping already existing Grafana user")
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
