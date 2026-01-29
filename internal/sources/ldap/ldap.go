// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package ldap

import (
	"crypto/tls"
	"log/slog"
	"os"
	"strconv"

	goldap "github.com/go-ldap/ldap/v3"

	"github.com/skuethe/grafana-oss-team-sync/internal/sources/sourcetypes"
)

type search struct {
	connection *goldap.Conn
	baseDN     string
	filter     string
	attributes []string
}

func (s *search) perform() (*goldap.SearchResult, error) {
	searchRequest := goldap.NewSearchRequest(
		s.baseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		s.filter,
		s.attributes,
		nil,
	)

	sr, err := s.connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) < 1 {
		slog.Warn("search result is empty")
	}

	return sr, nil
}

func New() *sourcetypes.SourcePlugin {
	ldapLog := slog.With(slog.String("package", "ldap"))
	ldapLog.Info("initializing LDAP")

	baseDN := os.Getenv("LDAP_BASE_DN")
	bindDN := os.Getenv("LDAP_BIND_DN")
	bindPassword := os.Getenv("LDAP_BIND_PASSWORD")
	uri := os.Getenv("LDAP_URI")
	groupFilter := os.Getenv("LDAP_GROUP_FILTER")
	userFilter := os.Getenv("LDAP_USER_FILTER")

	insecureSkipVerify, err := strconv.ParseBool(os.Getenv("LDAP_INSECURE_SKIP_VERIFY"))
	if err != nil {
		ldapLog.Error("could not convert variable LDAP_INSECURE_SKIP_VERIFY to bool")
	}

	conn, err := goldap.DialURL(uri)
	if err != nil {
		ldapLog.Error("unable to connect to LDAP URI")
		panic(err)
	}
	// defer conn.Close()

	// Reconnect with TLS
	if insecureSkipVerify {
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			ldapLog.Error("unable to connect via TLS")
			panic(err)
		}
	}

	// First bind with a read only user
	err = conn.Bind(bindDN, bindPassword)
	if err != nil {
		ldapLog.Error("unable to authenticate with bind user and password")
		panic(err)
	}

	return &sourcetypes.SourcePlugin{
		LDAP: &sourcetypes.LDAPClient{
			Connection:  conn,
			BaseDN:      baseDN,
			GroupFilter: groupFilter,
			UserFilter:  userFilter,
		},
	}
}
