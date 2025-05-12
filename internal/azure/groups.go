package azure

import (
	"log/slog"
	// "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// type groups struct {
// 	client *graph.GraphServiceClient
// 	log    slog.Logger
// }

// func (a *groups) doesGroupExist() bool {
// 	// sdkgroups.GroupItemRequestBuilderGetQueryParameters
// 	_, err := a.client.Groups().Get()
// 	return err == nil
// }

func (a *AzureInstance) processGroups() {
	groupsLog := slog.With(slog.String("package", "azure.groups"))
	groupsLog.Info("Initializing Azure Groups")

	countFound := 0
	countIgnored := 0

	// add logic

	groupsLog.Info(
		"Finished Azure Groups",
		slog.Group("stats",
			slog.Int("found", countFound),
			slog.Int("ignored", countIgnored),
		),
	)

}
