package migrator

import (
	"fmt"
	"strings"
	"github.com/go-xorm/xorm"
)

type Dialect interface {
	DriverName() string
	Quote(string) string
	AndStr() string
	AutoIncrStr() string
	OrStr() string
	EqStr() string
	ShowCreateNull() bool
	SqlType(col *Column) string
	SupportEngine() bool
	LikeStr() string
	Default(col *Column) string
	BooleanStr(bool) string
	DateTimeFunc(string) string
	CreateIndexSql(tableName string, index *Index) string
	CreateTableSql(table *Table) string
	AddColumnSql(tableName string, col *Column) string
	CopyTableData(sourceTable string, targetTable string, sourceCols []string, targetCols []string) string
	DropTable(tableName string) string
	DropIndexSql(tableName string, index *Index) string
	TableCheckSql(tableName string) (string, []interface{})
	RenameTable(oldName string, newName string) string
	UpdateTableSql(tableName string, columns []*Column) string
	ColString(*Column) string
	ColStringNoPk(*Column) string
	Limit(limit int64) string
	LimitOffset(limit int64, offset int64) string
	PreInsertId(table string, sess *xorm.Session) error
	PostInsertId(table string, sess *xorm.Session) error
	CleanDB() error
	NoOpSql() string
	IsUniqueConstraintViolation(err error) bool
}

func NewDialect(engine *xorm.Engine) Dialect {
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := engine.DriverName()
	switch name {
	case MYSQL:
		return NewMysqlDialect(engine)
	case SQLITE:
		return NewSqlite3Dialect(engine)
	case POSTGRES:
		return NewPostgresDialect(engine)
	}
	panic("Unsupported database type: " + name)
}

type BaseDialect struct {
	dialect		Dialect
	engine		*xorm.Engine
	driverName	string
}

func (d *BaseDialect) DriverName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.driverName
}
func (b *BaseDialect) ShowCreateNull() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (b *BaseDialect) AndStr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "AND"
}
func (b *BaseDialect) LikeStr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "LIKE"
}
func (b *BaseDialect) OrStr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "OR"
}
func (b *BaseDialect) EqStr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "="
}
func (b *BaseDialect) Default(col *Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return col.Default
}
func (db *BaseDialect) DateTimeFunc(value string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return value
}
func (b *BaseDialect) CreateTableSql(table *Table) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sql := "CREATE TABLE IF NOT EXISTS "
	sql += b.dialect.Quote(table.Name) + " (\n"
	pkList := table.PrimaryKeys
	for _, col := range table.Columns {
		if col.IsPrimaryKey && len(pkList) == 1 {
			sql += col.String(b.dialect)
		} else {
			sql += col.StringNoPk(b.dialect)
		}
		sql = strings.TrimSpace(sql)
		sql += "\n, "
	}
	if len(pkList) > 1 {
		quotedCols := []string{}
		for _, col := range pkList {
			quotedCols = append(quotedCols, b.dialect.Quote(col))
		}
		sql += "PRIMARY KEY ( " + strings.Join(quotedCols, ",") + " ), "
	}
	sql = sql[:len(sql)-2] + ")"
	if b.dialect.SupportEngine() {
		sql += " ENGINE=InnoDB DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci"
	}
	sql += ";"
	return sql
}
func (db *BaseDialect) AddColumnSql(tableName string, col *Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("alter table %s ADD COLUMN %s", db.dialect.Quote(tableName), col.StringNoPk(db.dialect))
}
func (db *BaseDialect) CreateIndexSql(tableName string, index *Index) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quote := db.dialect.Quote
	var unique string
	if index.Type == UniqueIndex {
		unique = " UNIQUE"
	}
	idxName := index.XName(tableName)
	quotedCols := []string{}
	for _, col := range index.Cols {
		quotedCols = append(quotedCols, db.dialect.Quote(col))
	}
	return fmt.Sprintf("CREATE%s INDEX %v ON %v (%v);", unique, quote(idxName), quote(tableName), strings.Join(quotedCols, ","))
}
func (db *BaseDialect) QuoteColList(cols []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var sourceColsSql = ""
	for _, col := range cols {
		sourceColsSql += db.dialect.Quote(col)
		sourceColsSql += "\n, "
	}
	return strings.TrimSuffix(sourceColsSql, "\n, ")
}
func (db *BaseDialect) CopyTableData(sourceTable string, targetTable string, sourceCols []string, targetCols []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sourceColsSql := db.QuoteColList(sourceCols)
	targetColsSql := db.QuoteColList(targetCols)
	quote := db.dialect.Quote
	return fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", quote(targetTable), targetColsSql, sourceColsSql, quote(sourceTable))
}
func (db *BaseDialect) DropTable(tableName string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quote := db.dialect.Quote
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", quote(tableName))
}
func (db *BaseDialect) RenameTable(oldName string, newName string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quote := db.dialect.Quote
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", quote(oldName), quote(newName))
}
func (db *BaseDialect) DropIndexSql(tableName string, index *Index) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quote := db.dialect.Quote
	name := index.XName(tableName)
	return fmt.Sprintf("DROP INDEX %v ON %s", quote(name), quote(tableName))
}
func (db *BaseDialect) UpdateTableSql(tableName string, columns []*Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "-- NOT REQUIRED"
}
func (db *BaseDialect) ColString(col *Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sql := db.dialect.Quote(col.Name) + " "
	sql += db.dialect.SqlType(col) + " "
	if col.IsPrimaryKey {
		sql += "PRIMARY KEY "
		if col.IsAutoIncrement {
			sql += db.dialect.AutoIncrStr() + " "
		}
	}
	if db.dialect.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}
	if col.Default != "" {
		sql += "DEFAULT " + db.dialect.Default(col) + " "
	}
	return sql
}
func (db *BaseDialect) ColStringNoPk(col *Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sql := db.dialect.Quote(col.Name) + " "
	sql += db.dialect.SqlType(col) + " "
	if db.dialect.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}
	if col.Default != "" {
		sql += "DEFAULT " + db.dialect.Default(col) + " "
	}
	return sql
}
func (db *BaseDialect) Limit(limit int64) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf(" LIMIT %d", limit)
}
func (db *BaseDialect) LimitOffset(limit int64, offset int64) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
func (db *BaseDialect) PreInsertId(table string, sess *xorm.Session) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (db *BaseDialect) PostInsertId(table string, sess *xorm.Session) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (db *BaseDialect) CleanDB() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (db *BaseDialect) NoOpSql() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "SELECT 0;"
}
