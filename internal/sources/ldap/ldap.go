// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package ldap

import (
	"log/slog"

	"github.com/skuethe/grafana-oss-team-sync/internal/sources/sourcetypes"
)

func New() *sourcetypes.SourcePlugin {
	ldapLog := slog.With(slog.String("package", "ldap"))
	ldapLog.Info("initializing LDAP")

	// clientId := os.Getenv("CLIENT_ID")
	// tenantId := os.Getenv("TENANT_ID")
	// clientSecret := os.Getenv("CLIENT_SECRET")

	// credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	// if err != nil {
	// 	ldapLog.Error("unable to create secret credential for msgraph client")
	// 	panic(err)
	// }

	// client, err := graph.NewGraphServiceClientWithCredentials(credential, []string{"https://graph.microsoft.com/.default"})
	// if err != nil {
	// 	ldapLog.Error("unable to create msgraph client")
	// 	panic(err)
	// }

	return &sourcetypes.SourcePlugin{
		LDAP: client,
	}
}
