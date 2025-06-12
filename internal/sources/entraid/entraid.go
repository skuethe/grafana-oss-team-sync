package entraid

import (
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/skuethe/grafana-oss-team-sync/internal/plugin"
)

func New() *plugin.SourceInstance {
	entraidLog := slog.With(slog.String("package", "entraid"))
	entraidLog.Info("initializing EntraID")

	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		entraidLog.Error("unable to create secret credential for msgraph client")
		panic(err)
	}

	client, err := graph.NewGraphServiceClientWithCredentials(
		credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		entraidLog.Error("unable to create msgraph client")
		panic(err)
	}

	return &plugin.SourceInstance{
		EntraID: client,
	}
}
