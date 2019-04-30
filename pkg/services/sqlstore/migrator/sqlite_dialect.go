package migrator

import (
	"fmt"
	"github.com/go-xorm/xorm"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct{ BaseDialect }

func NewSqlite3Dialect(engine *xorm.Engine) *Sqlite3 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d := Sqlite3{}
	d.BaseDialect.dialect = &d
	d.BaseDialect.engine = engine
	d.BaseDialect.driverName = SQLITE
	return &d
}
func (db *Sqlite3) SupportEngine() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (db *Sqlite3) Quote(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "`" + name + "`"
}
func (db *Sqlite3) AutoIncrStr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "AUTOINCREMENT"
}
func (db *Sqlite3) BooleanStr(value bool) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value {
		return "1"
	}
	return "0"
}
func (db *Sqlite3) DateTimeFunc(value string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "datetime(" + value + ")"
}
func (db *Sqlite3) SqlType(c *Column) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch c.Type {
	case DB_Date, DB_DateTime, DB_TimeStamp, DB_Time:
		return DB_DateTime
	case DB_TimeStampz:
		return DB_Text
	case DB_Char, DB_Varchar, DB_NVarchar, DB_TinyText, DB_Text, DB_MediumText, DB_LongText:
		return DB_Text
	case DB_Bit, DB_TinyInt, DB_SmallInt, DB_MediumInt, DB_Int, DB_Integer, DB_BigInt, DB_Bool:
		return DB_Integer
	case DB_Float, DB_Double, DB_Real:
		return DB_Real
	case DB_Decimal, DB_Numeric:
		return DB_Numeric
	case DB_TinyBlob, DB_Blob, DB_MediumBlob, DB_LongBlob, DB_Bytea, DB_Binary, DB_VarBinary:
		return DB_Blob
	case DB_Serial, DB_BigSerial:
		c.IsPrimaryKey = true
		c.IsAutoIncrement = true
		c.Nullable = false
		return DB_Integer
	default:
		return c.Type
	}
}
func (db *Sqlite3) TableCheckSql(tableName string) (string, []interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
func (db *Sqlite3) DropIndexSql(tableName string, index *Index) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quote := db.Quote
	idxName := index.XName(tableName)
	return fmt.Sprintf("DROP INDEX %v", quote(idxName))
}
func (db *Sqlite3) CleanDB() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (db *Sqlite3) IsUniqueConstraintViolation(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if driverErr, ok := err.(sqlite3.Error); ok {
		if driverErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true
		}
	}
	return false
}
