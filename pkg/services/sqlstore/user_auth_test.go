package sqlstore

import (
	"context"
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	m "github.com/grafana/grafana/pkg/models"
)

func TestUserAuth(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	InitTestDB(t)
	Convey("Given 5 users", t, func() {
		var err error
		var cmd *m.CreateUserCommand
		for i := 0; i < 5; i++ {
			cmd = &m.CreateUserCommand{Email: fmt.Sprint("user", i, "@test.com"), Name: fmt.Sprint("user", i), Login: fmt.Sprint("loginuser", i)}
			err = CreateUser(context.Background(), cmd)
			So(err, ShouldBeNil)
		}
		Reset(func() {
			_, err := x.Exec("DELETE FROM org_user WHERE 1=1")
			So(err, ShouldBeNil)
			_, err = x.Exec("DELETE FROM org WHERE 1=1")
			So(err, ShouldBeNil)
			_, err = x.Exec("DELETE FROM " + dialect.Quote("user") + " WHERE 1=1")
			So(err, ShouldBeNil)
			_, err = x.Exec("DELETE FROM user_auth WHERE 1=1")
			So(err, ShouldBeNil)
		})
		Convey("Can find existing user", func() {
			login := "loginuser0"
			query := &m.GetUserByAuthInfoQuery{Login: login}
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Login, ShouldEqual, login)
			id := query.Result.Id
			query = &m.GetUserByAuthInfoQuery{UserId: id}
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Id, ShouldEqual, id)
			email := "user1@test.com"
			query = &m.GetUserByAuthInfoQuery{Email: email}
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Email, ShouldEqual, email)
			email = "nonexistent@test.com"
			query = &m.GetUserByAuthInfoQuery{Email: email}
			err = GetUserByAuthInfo(query)
			So(err, ShouldEqual, m.ErrUserNotFound)
			So(query.Result, ShouldBeNil)
		})
		Convey("Can set & locate by AuthModule and AuthId", func() {
			query := &m.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			err = GetUserByAuthInfo(query)
			So(err, ShouldEqual, m.ErrUserNotFound)
			So(query.Result, ShouldBeNil)
			login := "loginuser0"
			query.Login = login
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Login, ShouldEqual, login)
			query = &m.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Login, ShouldEqual, login)
			id := query.Result.Id
			query.UserId = id + 1
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Login, ShouldEqual, "loginuser1")
			query = &m.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			err = GetUserByAuthInfo(query)
			So(err, ShouldBeNil)
			So(query.Result.Login, ShouldEqual, "loginuser1")
			_, err = x.Exec("DELETE FROM "+dialect.Quote("user")+" WHERE id=?", query.Result.Id)
			So(err, ShouldBeNil)
			query = &m.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			err = GetUserByAuthInfo(query)
			So(err, ShouldEqual, m.ErrUserNotFound)
			So(query.Result, ShouldBeNil)
		})
	})
}
