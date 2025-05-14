package azure

import (
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

func new() (azureInstance, error) {
	clientLog := slog.With(slog.String("package", "azure.client"))
	clientLog.Info("Initializing Azure Client")

	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return azureInstance{}, err
	}

	client, err := graph.NewGraphServiceClientWithCredentials(
		credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return azureInstance{}, err
	}

	return azureInstance{
		api: client,
	}, nil
}
