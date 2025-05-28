package plugin

import (
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type SourceInstance struct {
	EntraID *msgraph.GraphServiceClient
}
