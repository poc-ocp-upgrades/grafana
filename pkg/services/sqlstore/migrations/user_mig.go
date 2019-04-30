package migrations

import (
	"fmt"
	"github.com/go-xorm/xorm"
	. "github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util"
)

func addUserMigrations(mg *Migrator) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	userV1 := Table{Name: "user", Columns: []*Column{{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true}, {Name: "version", Type: DB_Int, Nullable: false}, {Name: "login", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "email", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "name", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "password", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "salt", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "rands", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "company", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "account_id", Type: DB_BigInt, Nullable: false}, {Name: "is_admin", Type: DB_Bool, Nullable: false}, {Name: "created", Type: DB_DateTime, Nullable: false}, {Name: "updated", Type: DB_DateTime, Nullable: false}}, Indices: []*Index{{Cols: []string{"login"}, Type: UniqueIndex}, {Cols: []string{"email"}, Type: UniqueIndex}}}
	mg.AddMigration("create user table", NewAddTableMigration(userV1))
	mg.AddMigration("add unique index user.login", NewAddIndexMigration(userV1, userV1.Indices[0]))
	mg.AddMigration("add unique index user.email", NewAddIndexMigration(userV1, userV1.Indices[1]))
	addDropAllIndicesMigrations(mg, "v1", userV1)
	addTableRenameMigration(mg, "user", "user_v1", "v1")
	userV2 := Table{Name: "user", Columns: []*Column{{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true}, {Name: "version", Type: DB_Int, Nullable: false}, {Name: "login", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "email", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "name", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "password", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "salt", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "rands", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "company", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "org_id", Type: DB_BigInt, Nullable: false}, {Name: "is_admin", Type: DB_Bool, Nullable: false}, {Name: "email_verified", Type: DB_Bool, Nullable: true}, {Name: "theme", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "created", Type: DB_DateTime, Nullable: false}, {Name: "updated", Type: DB_DateTime, Nullable: false}}, Indices: []*Index{{Cols: []string{"login"}, Type: UniqueIndex}, {Cols: []string{"email"}, Type: UniqueIndex}}}
	mg.AddMigration("create user table v2", NewAddTableMigration(userV2))
	addTableIndicesMigrations(mg, "v2", userV2)
	mg.AddMigration("copy data_source v1 to v2", NewCopyTableDataMigration("user", "user_v1", map[string]string{"id": "id", "version": "version", "login": "login", "email": "email", "name": "name", "password": "password", "salt": "salt", "rands": "rands", "company": "company", "org_id": "account_id", "is_admin": "is_admin", "created": "created", "updated": "updated"}))
	mg.AddMigration("Drop old table user_v1", NewDropTableMigration("user_v1"))
	mg.AddMigration("Add column help_flags1 to user table", NewAddColumnMigration(userV2, &Column{Name: "help_flags1", Type: DB_BigInt, Nullable: false, Default: "0"}))
	mg.AddMigration("Update user table charset", NewTableCharsetMigration("user", []*Column{{Name: "login", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "email", Type: DB_NVarchar, Length: 190, Nullable: false}, {Name: "name", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "password", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "salt", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "rands", Type: DB_NVarchar, Length: 50, Nullable: true}, {Name: "company", Type: DB_NVarchar, Length: 255, Nullable: true}, {Name: "theme", Type: DB_NVarchar, Length: 255, Nullable: true}}))
	mg.AddMigration("Add last_seen_at column to user", NewAddColumnMigration(userV2, &Column{Name: "last_seen_at", Type: DB_DateTime, Nullable: true}))
	mg.AddMigration("Add missing user data", &AddMissingUserSaltAndRandsMigration{})
}

type AddMissingUserSaltAndRandsMigration struct{ MigrationBase }

func (m *AddMissingUserSaltAndRandsMigration) Sql(dialect Dialect) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "code migration"
}

type TempUserDTO struct {
	Id	int64
	Login	string
}

func (m *AddMissingUserSaltAndRandsMigration) Exec(sess *xorm.Session, mg *Migrator) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	users := make([]*TempUserDTO, 0)
	err := sess.SQL(fmt.Sprintf("SELECT id, login from %s WHERE rands = ''", mg.Dialect.Quote("user"))).Find(&users)
	if err != nil {
		return err
	}
	for _, user := range users {
		_, err := sess.Exec("UPDATE "+mg.Dialect.Quote("user")+" SET salt = ?, rands = ? WHERE id = ?", util.GetRandomString(10), util.GetRandomString(10), user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
