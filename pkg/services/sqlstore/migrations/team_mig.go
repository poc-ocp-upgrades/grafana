package migrations

import . "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addTeamMigrations(mg *Migrator) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	teamV1 := Table{Name: "team", Columns: []*Column{{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true}, {Name: "name", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "org_id", Type: DB_BigInt}, {Name: "created", Type: DB_DateTime, Nullable: false}, {Name: "updated", Type: DB_DateTime, Nullable: false}}, Indices: []*Index{{Cols: []string{"org_id"}}, {Cols: []string{"org_id", "name"}, Type: UniqueIndex}}}
	mg.AddMigration("create team table", NewAddTableMigration(teamV1))
	mg.AddMigration("add index team.org_id", NewAddIndexMigration(teamV1, teamV1.Indices[0]))
	mg.AddMigration("add unique index team_org_id_name", NewAddIndexMigration(teamV1, teamV1.Indices[1]))
	teamMemberV1 := Table{Name: "team_member", Columns: []*Column{{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true}, {Name: "org_id", Type: DB_BigInt}, {Name: "team_id", Type: DB_BigInt}, {Name: "user_id", Type: DB_BigInt}, {Name: "created", Type: DB_DateTime, Nullable: false}, {Name: "updated", Type: DB_DateTime, Nullable: false}}, Indices: []*Index{{Cols: []string{"org_id"}}, {Cols: []string{"org_id", "team_id", "user_id"}, Type: UniqueIndex}}}
	mg.AddMigration("create team member table", NewAddTableMigration(teamMemberV1))
	mg.AddMigration("add index team_member.org_id", NewAddIndexMigration(teamMemberV1, teamMemberV1.Indices[0]))
	mg.AddMigration("add unique index team_member_org_id_team_id_user_id", NewAddIndexMigration(teamMemberV1, teamMemberV1.Indices[1]))
	mg.AddMigration("Add column email to team table", NewAddColumnMigration(teamV1, &Column{Name: "email", Type: DB_NVarchar, Nullable: true, Length: 190}))
	mg.AddMigration("Add column external to team_member table", NewAddColumnMigration(teamMemberV1, &Column{Name: "external", Type: DB_Bool, Nullable: true}))
}
