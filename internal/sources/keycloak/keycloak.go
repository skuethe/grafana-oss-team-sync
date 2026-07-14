package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
	kyClient "github.com/skuethe/grafana-oss-team-sync/internal/sources/keycloak/client"
)

type KeycloakParam struct {
	realm      string
	username   string
	password   string
	clientID   string
	url        string
	groupRegex *regexp.Regexp
}

type Keycloak struct {
	param  KeycloakParam
	client *kyClient.AdminClient
	logger *slog.Logger
}

func New(ctx context.Context) (*Keycloak, error) {
	keycloakLog := slog.With(slog.String("package", "keycloak"))
	keycloakLog.Info("initializing Keycloak")

	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	addr := os.Getenv("KEYCLOAK_URL")
	password := os.Getenv("KEYCLOAK_ADMIN_PASSWORD")
	username := os.Getenv("KEYCLOAK_ADMIN_USERNAME")
	realm := os.Getenv("KEYCLOAK_REALM")

	keycloakRegex, err := regexp.Compile(os.Getenv("KEYCLOAK_GROUP_REGEX"))
	if err != nil {
		return nil, fmt.Errorf("invalid keycloak groups regex: %w", err)
	}

	return &Keycloak{
		logger: keycloakLog,
		client: kyClient.NewAdminClient(addr, realm, username, password, clientID),
		param: KeycloakParam{
			realm:      realm,
			username:   username,
			password:   password,
			clientID:   clientID,
			url:        addr,
			groupRegex: keycloakRegex,
		},
	}, nil
}

// ProcessGroup recursively fetches all groups and filters by regex
func (k *Keycloak) ProcessGroups(ctx context.Context) *grafana.Teams {
	slog.Info("Fetching groups recursively")

	result := &grafana.Teams{}

	allGroups, err := k.client.BFSGroupTraversal(ctx)
	if err != nil {
		slog.Error("failed to get keycloak groups", "error", err)
		return nil
	}

	for key, v := range allGroups {
		if !k.param.groupRegex.MatchString(v) {
			continue
		}

		member, err := k.client.GetGroupMember(ctx, key)
		if err != nil {
			slog.Error("failed to get keycloak group members", "error", err)
			return nil
		}

		user := &grafana.Users{}
		for _, u := range member {
			*user = append(*user, grafana.User{
				Login: u.Username,
				Name:  u.Username,
				Email: u.Email,
			})
		}

		*result = append(*result, grafana.Team{
			Users: user,
			Parameter: &grafana.TeamParameter{
				Name: v,
			},
		})

	}

	return result
}
