package grafana

import (
	"log/slog"
	"os"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

type GrafanaInstance struct {
	api *client.GrafanaHTTPAPI
}

func Start(userList *[]models.AdminCreateUserForm, teamList *[]models.CreateTeamCommand) {
	grafanaLog := slog.With(slog.String("package", "grafana"))
	grafanaLog.Info("Initializing Grafana")

	instance := GrafanaInstance{
		api: initClient(),
	}

	grafanaLog.Info("Validating Grafana Health")
	_, err := instance.api.Health.GetHealth()
	if err != nil {
		grafanaLog.Error("Grafana instance is not healthy", "error", err)
		os.Exit(1)
	}

	instance.ProcessUsers(userList)
	instance.ProcessTeams(teamList)
	instance.ProcessFolders()
}
