package models

import (
	"time"
)

type LoginAttempt struct {
	Id			int64
	Username	string
	IpAddress	string
	Created		int64
}
type CreateLoginAttemptCommand struct {
	Username	string
	IpAddress	string
	Result		LoginAttempt
}
type DeleteOldLoginAttemptsCommand struct {
	OlderThan	time.Time
	DeletedRows	int64
}
type GetUserLoginAttemptCountQuery struct {
	Username	string
	Since		time.Time
	Result		int64
}
