package azure

import (
	"log/slog"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
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

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		clientLog.Error("Could not create Azure auth provider", "error", err)
		os.Exit(1)
	}

	// Create a request adapter using the auth provider
	adapter, err := graph.NewGraphRequestAdapter(authProvider)
	if err != nil {
		clientLog.Error("Could not create Azure request adapter", "error", err)
		os.Exit(1)
	}

	// Create a Graph client using request adapter
	client := graph.NewGraphServiceClient(adapter)
	return client

	// graphClient, err := graph.NewGraphServiceClientWithCredentials(
	// 	credential, []string{"https://graph.microsoft.com/.default"})
	// if err != nil {
	// 	clientLog.Error("Could not create Azure graph client", "error", err)
	// 	os.Exit(1)
	// }

	// return graphClient
}
