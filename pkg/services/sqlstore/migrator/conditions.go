package migrator

type MigrationCondition interface {
	Sql(dialect Dialect) (string, []interface{})
}
type IfTableExistsCondition struct{ TableName string }

func (c *IfTableExistsCondition) Sql(dialect Dialect) (string, []interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return dialect.TableCheckSql(c.TableName)
}
