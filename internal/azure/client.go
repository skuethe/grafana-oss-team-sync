package azure

import (
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
)

func initClient() *graph.GraphServiceClient {
	clientLog := slog.With(slog.String("package", "azure.client"))
	clientLog.Info("Initializing Azure Client")

	clientId := os.Getenv("CLIENT_ID")
	tenantId := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		clientLog.Error("Could not create Azure credential client", "error", err)
		os.Exit(1)
	}

	graphClient, err := graph.NewGraphServiceClientWithCredentials(credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		clientLog.Error("Could not create Azure graph client", "error", err)
		os.Exit(1)
	}

	return graphClient
}
