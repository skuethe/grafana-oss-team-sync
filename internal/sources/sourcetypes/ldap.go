// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package sourcetypes

import "github.com/go-ldap/ldap/v3"

type LDAPClient struct {
	Connection  *ldap.Conn
	BaseDN      string
	GroupFilter string
	UserFilter  string
}
