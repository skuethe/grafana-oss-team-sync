// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package grafana

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/teams"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/skuethe/grafana-oss-team-sync/internal/config"
)

type TeamParameter models.CreateTeamCommand

type Team struct {
	Parameter *TeamParameter
	Users     *Users
	// OrgID is the Grafana organization to sync this team into.
	// 0 (unset) uses whichever org the configured Grafana credentials default to.
	OrgID int64
}

type Teams []Team

// api returns the Grafana API client scoped to this team's organization.
func (t *Team) api() *client.GrafanaHTTPAPI {
	if t.OrgID == 0 {
		return Instance.api
	}
	return Instance.api.WithOrgID(t.OrgID)
}

// ensureOrgMembership makes sure the given user/login is a member of this team's organization,
// which is required before it can be added as a member of a team in that organization.
func (t *Team) ensureOrgMembership(loginOrEmail string) error {
	if t.OrgID == 0 {
		return nil
	}

	orgUsers, err := Instance.api.Orgs.GetOrgUsers(t.OrgID)
	if err != nil {
		return err
	}
	for _, orgUser := range orgUsers.Payload {
		if strings.EqualFold(orgUser.Login, loginOrEmail) || strings.EqualFold(orgUser.Email, loginOrEmail) {
			return nil
		}
	}

	_, err = Instance.api.Orgs.AddOrgUser(t.OrgID, &models.AddOrgUserCommand{
		LoginOrEmail: loginOrEmail,
		Role:         "Viewer",
	})
	return err
}

func (t *Team) searchTeam() (*teams.SearchTeamsOK, error) {
	result, err := t.api().Teams.SearchTeams(&teams.SearchTeamsParams{
		Name: t.Parameter.Name,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Team) doesTeamExist() (bool, error) {
	result, err := t.searchTeam()
	if err != nil {
		return false, err
	}
	return len(result.Payload.Teams) == 1, nil
}

func (t *Team) getTeamID() (*int64, error) {
	result, err := t.searchTeam()
	if err != nil {
		return nil, err
	}
	return result.Payload.Teams[0].ID, nil
}

func (t *Team) createTeam() error {
	_, err := t.api().Teams.CreateTeam(&models.CreateTeamCommand{
		Name:  t.Parameter.Name,
		Email: t.Parameter.Email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *Team) addUsersToTeam() (*[]string, error) {
	teamID, tErr := t.getTeamID()
	if tErr != nil {
		return nil, tErr
	}

	adminMemberList := &[]string{}
	teamMemberList := &[]string{}

	if config.Instance.Features.AddLocalAdminToTeams {
		if err := t.ensureOrgMembership("admin@localhost"); err != nil {
			slog.Warn("could not add local admin to organization",
				slog.Int64("orgId", t.OrgID),
				slog.Any("error", err),
			)
		} else {
			*adminMemberList = append(*adminMemberList, "admin@localhost")
		}
	}

	for _, user := range *t.Users {
		if user.Email == "" {
			slog.Warn("skipping user with no mail",
				slog.Group("user",
					slog.String("login", user.Login),
					slog.String("name", user.Name),
				),
			)
			continue
		}
		if err := t.ensureOrgMembership(user.Email); err != nil {
			slog.Warn("could not add user to organization, skipping team membership",
				slog.String("user", user.Email),
				slog.Int64("orgId", t.OrgID),
				slog.Any("error", err),
			)
			continue
		}
		*teamMemberList = append(*teamMemberList, user.Email)
	}

	if _, err := t.api().Teams.SetTeamMemberships(strconv.FormatInt(*teamID, 10), &models.SetTeamMembershipsCommand{
		Admins:  *adminMemberList,
		Members: *teamMemberList,
	}); err != nil {
		return nil, fmt.Errorf("%w. Inputs: %v; %v", err, *adminMemberList, *teamMemberList)
	}

	return teamMemberList, nil
}

func (t *Teams) ProcessTeams() {
	teamsLog := slog.With(slog.String("package", "grafana.teams"))

	if len(*t) == 0 {
		teamsLog.Info("no teams to process, skipping")
	} else {
		teamsLog.Info("processing teams")

		countSkipped := 0
		countCreated := 0

		for _, team := range *t {

			teamLog := slog.With(
				slog.Group("team",
					slog.String("name", *team.Parameter.Name),
					slog.Int64("orgId", team.OrgID),
				),
			)

			exists, err := team.doesTeamExist()
			if err != nil {
				teamLog.Error("could not search for team",
					slog.Any("error", err),
				)
			} else {
				if exists {
					countSkipped++
					teamLog.Debug("skipping already existing team")
				} else {
					err := team.createTeam()
					if err != nil {
						teamLog.Error("could not create team",
							slog.Any("error", err),
						)
						continue
					} else {
						teamLog.Info("created team")
						countCreated++
					}
				}
				// Add users to team, if userSync feature is enabled
				if !config.Instance.Features.DisableUserSync {
					teamLog.Info("processing team members")

					if userList, err := team.addUsersToTeam(); err != nil {
						teamLog.Error("could not add users to team",
							slog.Any("error", err),
						)
					} else {
						teamLog.Debug("added users to team",
							slog.Any("list", *userList),
						)
					}

					teamLog.Info("finished processing team members")
				}
			}
		}

		teamsLog.Info("finished processing teams",
			slog.Group("teams",
				slog.Int("created", countCreated),
				slog.Int("existing", countSkipped),
			),
		)
	}
}
