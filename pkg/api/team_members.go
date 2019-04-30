package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

func GetTeamMembers(c *m.ReqContext) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	query := m.GetTeamMembersQuery{OrgId: c.OrgId, TeamId: c.ParamsInt64(":teamId")}
	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Failed to get Team Members", err)
	}
	for _, member := range query.Result {
		member.AvatarUrl = dtos.GetGravatarUrl(member.Email)
		member.Labels = []string{}
		if setting.IsEnterprise && setting.LdapEnabled && member.External {
			member.Labels = append(member.Labels, "LDAP")
		}
	}
	return JSON(200, query.Result)
}
func AddTeamMember(c *m.ReqContext, cmd m.AddTeamMemberCommand) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd.TeamId = c.ParamsInt64(":teamId")
	cmd.OrgId = c.OrgId
	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Team not found", nil)
		}
		if err == m.ErrTeamMemberAlreadyAdded {
			return Error(400, "User is already added to this team", nil)
		}
		return Error(500, "Failed to add Member to Team", err)
	}
	return JSON(200, &util.DynMap{"message": "Member added to Team"})
}
func RemoveTeamMember(c *m.ReqContext) Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := bus.Dispatch(&m.RemoveTeamMemberCommand{OrgId: c.OrgId, TeamId: c.ParamsInt64(":teamId"), UserId: c.ParamsInt64(":userId")}); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Team not found", nil)
		}
		if err == m.ErrTeamMemberNotFound {
			return Error(404, "Team member not found", nil)
		}
		return Error(500, "Failed to remove Member from Team", err)
	}
	return Success("Team Member removed")
}
