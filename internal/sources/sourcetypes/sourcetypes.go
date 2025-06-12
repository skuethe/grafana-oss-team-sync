package sourcetypes

import (
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
)

type SourcePlugin struct {
	EntraID *msgraph.GraphServiceClient
}
