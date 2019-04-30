package sqlstore

import (
	"time"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bus.AddHandler("sql", GetUserByAuthInfo)
	bus.AddHandler("sql", GetAuthInfo)
	bus.AddHandler("sql", SetAuthInfo)
	bus.AddHandler("sql", DeleteAuthInfo)
}
func GetUserByAuthInfo(query *m.GetUserByAuthInfoQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	user := &m.User{}
	has := false
	var err error
	authQuery := &m.GetAuthInfoQuery{}
	if query.AuthModule != "" && query.AuthId != "" {
		authQuery.AuthModule = query.AuthModule
		authQuery.AuthId = query.AuthId
		err = GetAuthInfo(authQuery)
		if err != m.ErrUserNotFound {
			if err != nil {
				return err
			}
			if query.UserId != 0 && query.UserId != authQuery.Result.UserId {
				err = DeleteAuthInfo(&m.DeleteAuthInfoCommand{UserAuth: authQuery.Result})
				if err != nil {
					sqlog.Error("Error removing user_auth entry", "error", err)
				}
				authQuery.Result = nil
			} else {
				has, err = x.Id(authQuery.Result.UserId).Get(user)
				if err != nil {
					return err
				}
				if !has {
					err = DeleteAuthInfo(&m.DeleteAuthInfoCommand{UserAuth: authQuery.Result})
					if err != nil {
						sqlog.Error("Error removing user_auth entry", "error", err)
					}
					authQuery.Result = nil
				}
			}
		}
	}
	if !has && query.UserId != 0 {
		has, err = x.Id(query.UserId).Get(user)
		if err != nil {
			return err
		}
	}
	if !has && query.Email != "" {
		user = &m.User{Email: query.Email}
		has, err = x.Get(user)
		if err != nil {
			return err
		}
	}
	if !has && query.Login != "" {
		user = &m.User{Login: query.Login}
		has, err = x.Get(user)
		if err != nil {
			return err
		}
	}
	if !has {
		return m.ErrUserNotFound
	}
	if authQuery.Result == nil && query.AuthModule != "" && query.AuthId != "" {
		cmd2 := &m.SetAuthInfoCommand{UserId: user.Id, AuthModule: query.AuthModule, AuthId: query.AuthId}
		if err := SetAuthInfo(cmd2); err != nil {
			return err
		}
	}
	query.Result = user
	return nil
}
func GetAuthInfo(query *m.GetAuthInfoQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	userAuth := &m.UserAuth{AuthModule: query.AuthModule, AuthId: query.AuthId}
	has, err := x.Get(userAuth)
	if err != nil {
		return err
	}
	if !has {
		return m.ErrUserNotFound
	}
	query.Result = userAuth
	return nil
}
func SetAuthInfo(cmd *m.SetAuthInfoCommand) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return inTransaction(func(sess *DBSession) error {
		authUser := &m.UserAuth{UserId: cmd.UserId, AuthModule: cmd.AuthModule, AuthId: cmd.AuthId, Created: time.Now()}
		_, err := sess.Insert(authUser)
		return err
	})
}
func DeleteAuthInfo(cmd *m.DeleteAuthInfoCommand) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return inTransaction(func(sess *DBSession) error {
		_, err := sess.Delete(cmd.UserAuth)
		return err
	})
}
