package sqlstore

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/annotations"
	"github.com/grafana/grafana/pkg/services/cache"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrations"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/services/sqlstore/sqlutil"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/grafana/grafana/pkg/tsdb/mssql"
	_ "github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	x	*xorm.Engine
	dialect	migrator.Dialect
	sqlog	log.Logger	= log.New("sqlstore")
)

const ContextSessionName = "db-session"

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registry.Register(&registry.Descriptor{Name: "SqlStore", Instance: &SqlStore{}, InitPriority: registry.High})
}

type SqlStore struct {
	Cfg		*setting.Cfg		`inject:""`
	Bus		bus.Bus			`inject:""`
	CacheService	*cache.CacheService	`inject:""`
	dbCfg		DatabaseConfig
	engine		*xorm.Engine
	log		log.Logger
	Dialect		migrator.Dialect
	skipEnsureAdmin	bool
}

func (ss *SqlStore) NewSession() *DBSession {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &DBSession{Session: ss.engine.NewSession()}
}
func (ss *SqlStore) WithDbSession(ctx context.Context, callback dbTransactionFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sess, err := startSession(ctx, ss.engine, false)
	if err != nil {
		return err
	}
	return callback(sess)
}
func (ss *SqlStore) WithTransactionalDbSession(ctx context.Context, callback dbTransactionFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ss.inTransactionWithRetryCtx(ctx, callback, 0)
}
func (ss *SqlStore) inTransactionWithRetryCtx(ctx context.Context, callback dbTransactionFunc, retry int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sess, err := startSession(ctx, ss.engine, true)
	if err != nil {
		return err
	}
	defer sess.Close()
	err = callback(sess)
	if sqlError, ok := err.(sqlite3.Error); ok && retry < 5 {
		if sqlError.Code == sqlite3.ErrLocked {
			sess.Rollback()
			time.Sleep(time.Millisecond * time.Duration(10))
			sqlog.Info("Database table locked, sleeping then retrying", "retry", retry)
			return ss.inTransactionWithRetryCtx(ctx, callback, retry+1)
		}
	}
	if err != nil {
		sess.Rollback()
		return err
	} else if err = sess.Commit(); err != nil {
		return err
	}
	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.Publish(e); err != nil {
				log.Error(3, "Failed to publish event after commit. error: %v", err)
			}
		}
	}
	return nil
}
func (ss *SqlStore) Init() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ss.log = log.New("sqlstore")
	ss.readConfig()
	engine, err := ss.getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}
	ss.engine = engine
	ss.Dialect = migrator.NewDialect(ss.engine)
	x = engine
	dialect = ss.Dialect
	migrator := migrator.NewMigrator(x)
	migrations.AddMigrations(migrator)
	for _, descriptor := range registry.GetServices() {
		sc, ok := descriptor.Instance.(registry.DatabaseMigrator)
		if ok {
			sc.AddMigration(migrator)
		}
	}
	if err := migrator.Start(); err != nil {
		return fmt.Errorf("Migration failed err: %v", err)
	}
	annotations.SetRepository(&SqlAnnotationRepo{})
	ss.Bus.SetTransactionManager(ss)
	ss.addUserQueryAndCommandHandlers()
	if ss.skipEnsureAdmin {
		return nil
	}
	return ss.ensureAdminUser()
}
func (ss *SqlStore) ensureAdminUser() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	systemUserCountQuery := m.GetSystemUserCountStatsQuery{}
	err := ss.InTransaction(context.Background(), func(ctx context.Context) error {
		err := bus.DispatchCtx(ctx, &systemUserCountQuery)
		if err != nil {
			return fmt.Errorf("Could not determine if admin user exists: %v", err)
		}
		if systemUserCountQuery.Result.Count > 0 {
			return nil
		}
		cmd := m.CreateUserCommand{}
		cmd.Login = setting.AdminUser
		cmd.Email = setting.AdminUser + "@localhost"
		cmd.Password = setting.AdminPassword
		cmd.IsAdmin = true
		if err := bus.DispatchCtx(ctx, &cmd); err != nil {
			return fmt.Errorf("Failed to create admin user: %v", err)
		}
		ss.log.Info("Created default admin", "user", setting.AdminUser)
		return nil
	})
	return err
}
func (ss *SqlStore) buildConnectionString() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cnnstr := ss.dbCfg.ConnectionString
	if cnnstr != "" {
		return cnnstr, nil
	}
	switch ss.dbCfg.Type {
	case migrator.MYSQL:
		protocol := "tcp"
		if strings.HasPrefix(ss.dbCfg.Host, "/") {
			protocol = "unix"
		}
		cnnstr = fmt.Sprintf("%s:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true", ss.dbCfg.User, ss.dbCfg.Pwd, protocol, ss.dbCfg.Host, ss.dbCfg.Name)
		if ss.dbCfg.SslMode == "true" || ss.dbCfg.SslMode == "skip-verify" {
			tlsCert, err := makeCert("custom", ss.dbCfg)
			if err != nil {
				return "", err
			}
			mysql.RegisterTLSConfig("custom", tlsCert)
			cnnstr += "&tls=custom"
		}
	case migrator.POSTGRES:
		var host, port = "127.0.0.1", "5432"
		fields := strings.Split(ss.dbCfg.Host, ":")
		if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
			host = fields[0]
		}
		if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
			port = fields[1]
		}
		if ss.dbCfg.Pwd == "" {
			ss.dbCfg.Pwd = "''"
		}
		if ss.dbCfg.User == "" {
			ss.dbCfg.User = "''"
		}
		cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", ss.dbCfg.User, ss.dbCfg.Pwd, host, port, ss.dbCfg.Name, ss.dbCfg.SslMode, ss.dbCfg.ClientCertPath, ss.dbCfg.ClientKeyPath, ss.dbCfg.CaCertPath)
	case migrator.SQLITE:
		if !filepath.IsAbs(ss.dbCfg.Path) {
			ss.dbCfg.Path = filepath.Join(ss.Cfg.DataPath, ss.dbCfg.Path)
		}
		os.MkdirAll(path.Dir(ss.dbCfg.Path), os.ModePerm)
		cnnstr = "file:" + ss.dbCfg.Path + "?cache=shared&mode=rwc"
	default:
		return "", fmt.Errorf("Unknown database type: %s", ss.dbCfg.Type)
	}
	return cnnstr, nil
}
func (ss *SqlStore) getEngine() (*xorm.Engine, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	connectionString, err := ss.buildConnectionString()
	if err != nil {
		return nil, err
	}
	sqlog.Info("Connecting to DB", "dbtype", ss.dbCfg.Type)
	engine, err := xorm.NewEngine(ss.dbCfg.Type, connectionString)
	if err != nil {
		return nil, err
	}
	engine.SetMaxOpenConns(ss.dbCfg.MaxOpenConn)
	engine.SetMaxIdleConns(ss.dbCfg.MaxIdleConn)
	engine.SetConnMaxLifetime(time.Second * time.Duration(ss.dbCfg.ConnMaxLifetime))
	debugSql := ss.Cfg.Raw.Section("database").Key("log_queries").MustBool(false)
	if !debugSql {
		engine.SetLogger(&xorm.DiscardLogger{})
	} else {
		engine.SetLogger(NewXormLogger(log.LvlInfo, log.New("sqlstore.xorm")))
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
	}
	return engine, nil
}
func (ss *SqlStore) readConfig() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sec := ss.Cfg.Raw.Section("database")
	cfgURL := sec.Key("url").String()
	if len(cfgURL) != 0 {
		dbURL, _ := url.Parse(cfgURL)
		ss.dbCfg.Type = dbURL.Scheme
		ss.dbCfg.Host = dbURL.Host
		pathSplit := strings.Split(dbURL.Path, "/")
		if len(pathSplit) > 1 {
			ss.dbCfg.Name = pathSplit[1]
		}
		userInfo := dbURL.User
		if userInfo != nil {
			ss.dbCfg.User = userInfo.Username()
			ss.dbCfg.Pwd, _ = userInfo.Password()
		}
	} else {
		ss.dbCfg.Type = sec.Key("type").String()
		ss.dbCfg.Host = sec.Key("host").String()
		ss.dbCfg.Name = sec.Key("name").String()
		ss.dbCfg.User = sec.Key("user").String()
		ss.dbCfg.ConnectionString = sec.Key("connection_string").String()
		ss.dbCfg.Pwd = sec.Key("password").String()
	}
	ss.dbCfg.MaxOpenConn = sec.Key("max_open_conn").MustInt(0)
	ss.dbCfg.MaxIdleConn = sec.Key("max_idle_conn").MustInt(2)
	ss.dbCfg.ConnMaxLifetime = sec.Key("conn_max_lifetime").MustInt(14400)
	ss.dbCfg.SslMode = sec.Key("ssl_mode").String()
	ss.dbCfg.CaCertPath = sec.Key("ca_cert_path").String()
	ss.dbCfg.ClientKeyPath = sec.Key("client_key_path").String()
	ss.dbCfg.ClientCertPath = sec.Key("client_cert_path").String()
	ss.dbCfg.ServerCertName = sec.Key("server_cert_name").String()
	ss.dbCfg.Path = sec.Key("path").MustString("data/grafana.db")
}
func InitTestDB(t *testing.T) *SqlStore {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Helper()
	sqlstore := &SqlStore{}
	sqlstore.skipEnsureAdmin = true
	sqlstore.Bus = bus.New()
	sqlstore.CacheService = cache.New(5*time.Minute, 10*time.Minute)
	dbType := migrator.SQLITE
	if db, present := os.LookupEnv("GRAFANA_TEST_DB"); present {
		dbType = db
	}
	sqlstore.Cfg = setting.NewCfg()
	sec, _ := sqlstore.Cfg.Raw.NewSection("database")
	sec.NewKey("type", dbType)
	switch dbType {
	case "mysql":
		sec.NewKey("connection_string", sqlutil.TestDB_Mysql.ConnStr)
	case "postgres":
		sec.NewKey("connection_string", sqlutil.TestDB_Postgres.ConnStr)
	default:
		sec.NewKey("connection_string", sqlutil.TestDB_Sqlite3.ConnStr)
	}
	engine, err := xorm.NewEngine(dbType, sec.Key("connection_string").String())
	if err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}
	sqlstore.Dialect = migrator.NewDialect(engine)
	dialect = sqlstore.Dialect
	if err := dialect.CleanDB(); err != nil {
		t.Fatalf("Failed to clean test db %v", err)
	}
	if err := sqlstore.Init(); err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}
	sqlstore.engine.DatabaseTZ = time.UTC
	sqlstore.engine.TZLocation = time.UTC
	return sqlstore
}
func IsTestDbMySql() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if db, present := os.LookupEnv("GRAFANA_TEST_DB"); present {
		return db == migrator.MYSQL
	}
	return false
}
func IsTestDbPostgres() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if db, present := os.LookupEnv("GRAFANA_TEST_DB"); present {
		return db == migrator.POSTGRES
	}
	return false
}

type DatabaseConfig struct {
	Type, Host, Name, User, Pwd, Path, SslMode	string
	CaCertPath					string
	ClientKeyPath					string
	ClientCertPath					string
	ServerCertName					string
	ConnectionString				string
	MaxOpenConn					int
	MaxIdleConn					int
	ConnMaxLifetime					int
}
