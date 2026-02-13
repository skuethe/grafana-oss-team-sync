// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package entraid

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/skuethe/grafana-oss-team-sync/internal/config"
	"github.com/skuethe/grafana-oss-team-sync/internal/grafana"
)

type EntraID struct {
	client   *graph.GraphServiceClient
	clientID string
	tenantID string
}

func New() *EntraID {
	entraidLog := slog.With(slog.String("package", "entraid"))
	entraidLog.Info("initializing EntraID")

	clientID := os.Getenv("CLIENT_ID")
	tenantID := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	credential, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		entraidLog.Error("unable to create secret credential for msgraph client")
		panic(err)
	}

	client, err := graph.NewGraphServiceClientWithCredentials(credential, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		entraidLog.Error("unable to create msgraph client")
		panic(err)
	}

	return &EntraID{
		client:   client,
		clientID: clientID,
		tenantID: tenantID,
	}
}

func (t *EntraID) ProcessGroups(ctx context.Context) *grafana.Teams {
	groupsLog := slog.With(slog.String("package", "entraid.groups"))
	groupsLog.Info("processing EntraID groups")

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	g := groups{
		client:       t.client,
		headers:      headers,
		grafanaTeams: &grafana.Teams{},
	}

	gr, err := g.getInitialGroupRequest()
	if err != nil {
		groupsLog.Error("could not get initial group result from EntraID")
		panic(err)
	}

	countFound := (*gr).GetOdataCount()

	for {
		// Handle group result
		g.processGroupResult(gr)

		// Handle possible pagination
		nextPageURL := (*gr).GetOdataNextLink()
		if nextPageURL != nil {
			groupsLog.Debug("processing paginated group result")
			gr, err = g.handleGroupPagination(nextPageURL)
			if err != nil {
				groupsLog.Error("could not get paged group result from EntraID")
				panic(err)
			}
		} else {
			break
		}
	}

	if len(config.Instance.Teams) > 0 {
		groupsLog.Warn("could not find the following groups in EntraID", "skipped", strings.Join(config.Instance.Teams, ","))
	}

	groupsLog.Info("finished processing EntraID groups",
		slog.Group("groups",
			slog.Int64("found", *countFound),
			slog.Int("skipped", len(config.Instance.Teams)),
		),
	)

	return g.grafanaTeams
}
