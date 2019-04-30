package sqlstore

import (
	"strings"
	m "github.com/grafana/grafana/pkg/models"
)

type SearchBuilder struct {
	SqlBuilder
	tags			[]string
	isStarred		bool
	limit			int
	signedInUser		*m.SignedInUser
	whereDashboardIdsIn	[]int64
	whereTitle		string
	whereTypeFolder		bool
	whereTypeDash		bool
	whereFolderIds		[]int64
	permission		m.PermissionType
}

func NewSearchBuilder(signedInUser *m.SignedInUser, limit int, permission m.PermissionType) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	searchBuilder := &SearchBuilder{signedInUser: signedInUser, limit: limit, permission: permission}
	return searchBuilder
}
func (sb *SearchBuilder) WithTags(tags []string) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(tags) > 0 {
		sb.tags = tags
	}
	return sb
}
func (sb *SearchBuilder) IsStarred() *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.isStarred = true
	return sb
}
func (sb *SearchBuilder) WithDashboardIdsIn(ids []int64) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(ids) > 0 {
		sb.whereDashboardIdsIn = ids
	}
	return sb
}
func (sb *SearchBuilder) WithTitle(title string) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.whereTitle = title
	return sb
}
func (sb *SearchBuilder) WithType(queryType string) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(queryType) > 0 && queryType == "dash-folder" {
		sb.whereTypeFolder = true
	}
	if len(queryType) > 0 && queryType == "dash-db" {
		sb.whereTypeDash = true
	}
	return sb
}
func (sb *SearchBuilder) WithFolderIds(folderIds []int64) *SearchBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.whereFolderIds = folderIds
	return sb
}
func (sb *SearchBuilder) ToSql() (string, []interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.params = make([]interface{}, 0)
	sb.buildSelect()
	if len(sb.tags) > 0 {
		sb.buildTagQuery()
	} else {
		sb.buildMainQuery()
	}
	sb.sql.WriteString(`
		LEFT OUTER JOIN dashboard folder on folder.id = dashboard.folder_id
		LEFT OUTER JOIN dashboard_tag on dashboard.id = dashboard_tag.dashboard_id`)
	sb.sql.WriteString(" ORDER BY dashboard.title ASC" + dialect.Limit(5000))
	return sb.sql.String(), sb.params
}
func (sb *SearchBuilder) buildSelect() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.sql.WriteString(`SELECT
			dashboard.id,
			dashboard.uid,
			dashboard.title,
			dashboard.slug,
			dashboard_tag.term,
			dashboard.is_folder,
			dashboard.folder_id,
			folder.uid as folder_uid,
			folder.slug as folder_slug,
			folder.title as folder_title
		FROM `)
}
func (sb *SearchBuilder) buildTagQuery() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.sql.WriteString(`(
	SELECT
		dashboard.id FROM dashboard
		LEFT OUTER JOIN dashboard_tag ON dashboard_tag.dashboard_id = dashboard.id
	`)
	if sb.isStarred {
		sb.sql.WriteString(" INNER JOIN star on star.dashboard_id = dashboard.id")
	}
	sb.sql.WriteString(` WHERE dashboard_tag.term IN (?` + strings.Repeat(",?", len(sb.tags)-1) + `) AND `)
	for _, tag := range sb.tags {
		sb.params = append(sb.params, tag)
	}
	sb.buildSearchWhereClause()
	sb.sql.WriteString(`
		GROUP BY dashboard.id HAVING COUNT(dashboard.id) >= ?
		ORDER BY dashboard.id` + dialect.Limit(int64(sb.limit)) + `) as ids
		INNER JOIN dashboard on ids.id = dashboard.id
	`)
	sb.params = append(sb.params, len(sb.tags))
}
func (sb *SearchBuilder) buildMainQuery() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.sql.WriteString(`( SELECT dashboard.id FROM dashboard `)
	if sb.isStarred {
		sb.sql.WriteString(" INNER JOIN star on star.dashboard_id = dashboard.id")
	}
	sb.sql.WriteString(` WHERE `)
	sb.buildSearchWhereClause()
	sb.sql.WriteString(` ORDER BY dashboard.title` + dialect.Limit(int64(sb.limit)) + `) as ids INNER JOIN dashboard on ids.id = dashboard.id `)
}
func (sb *SearchBuilder) buildSearchWhereClause() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sb.sql.WriteString(` dashboard.org_id=?`)
	sb.params = append(sb.params, sb.signedInUser.OrgId)
	if sb.isStarred {
		sb.sql.WriteString(` AND star.user_id=?`)
		sb.params = append(sb.params, sb.signedInUser.UserId)
	}
	if len(sb.whereDashboardIdsIn) > 0 {
		sb.sql.WriteString(` AND dashboard.id IN (?` + strings.Repeat(",?", len(sb.whereDashboardIdsIn)-1) + `)`)
		for _, dashboardId := range sb.whereDashboardIdsIn {
			sb.params = append(sb.params, dashboardId)
		}
	}
	sb.writeDashboardPermissionFilter(sb.signedInUser, sb.permission)
	if len(sb.whereTitle) > 0 {
		sb.sql.WriteString(" AND dashboard.title " + dialect.LikeStr() + " ?")
		sb.params = append(sb.params, "%"+sb.whereTitle+"%")
	}
	if sb.whereTypeFolder {
		sb.sql.WriteString(" AND dashboard.is_folder = " + dialect.BooleanStr(true))
	}
	if sb.whereTypeDash {
		sb.sql.WriteString(" AND dashboard.is_folder = " + dialect.BooleanStr(false))
	}
	if len(sb.whereFolderIds) > 0 {
		sb.sql.WriteString(` AND dashboard.folder_id IN (?` + strings.Repeat(",?", len(sb.whereFolderIds)-1) + `) `)
		for _, id := range sb.whereFolderIds {
			sb.params = append(sb.params, id)
		}
	}
}
