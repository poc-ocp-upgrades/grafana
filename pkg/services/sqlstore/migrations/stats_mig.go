package migrations

import . "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addTestDataMigrations(mg *Migrator) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testData := Table{Name: "test_data", Columns: []*Column{{Name: "id", Type: DB_Int, IsPrimaryKey: true, IsAutoIncrement: true}, {Name: "metric1", Type: DB_Varchar, Length: 20, Nullable: true}, {Name: "metric2", Type: DB_NVarchar, Length: 150, Nullable: true}, {Name: "value_big_int", Type: DB_BigInt, Nullable: true}, {Name: "value_double", Type: DB_Double, Nullable: true}, {Name: "value_float", Type: DB_Float, Nullable: true}, {Name: "value_int", Type: DB_Int, Nullable: true}, {Name: "time_epoch", Type: DB_BigInt, Nullable: false}, {Name: "time_date_time", Type: DB_DateTime, Nullable: false}, {Name: "time_time_stamp", Type: DB_TimeStamp, Nullable: false}}}
	mg.AddMigration("create test_data table", NewAddTableMigration(testData))
}
