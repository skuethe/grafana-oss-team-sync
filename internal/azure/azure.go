package azure

import (
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

type AzureInstance struct {
	api *graph.GraphServiceClient
}

func New() *AzureInstance {
	azureLog := slog.With(slog.String("package", "azure"))
	azureLog.Info("initializing Azure")

	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		azureLog.Error("unable to create secret credential for client", "error", err)
		os.Exit(1)
	}

	client, err := graph.NewGraphServiceClientWithCredentials(
		credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		azureLog.Error("unable to create client", "error", err)
		os.Exit(1)
	}

	return &AzureInstance{
		api: client,
	}
}

func Start() {
	instance := New()

	instance.processGroups()

	// instance.processUsers(userList)
}
