package models

import (
	"errors"
	"time"
)

var (
	ErrTeamMemberAlreadyAdded = errors.New("User is already added to this team")
)

type TeamMember struct {
	Id		int64
	OrgId		int64
	TeamId		int64
	UserId		int64
	External	bool
	Created		time.Time
	Updated		time.Time
}
type AddTeamMemberCommand struct {
	UserId		int64	`json:"userId" binding:"Required"`
	OrgId		int64	`json:"-"`
	TeamId		int64	`json:"-"`
	External	bool	`json:"-"`
}
type RemoveTeamMemberCommand struct {
	OrgId	int64	`json:"-"`
	UserId	int64
	TeamId	int64
}
type GetTeamMembersQuery struct {
	OrgId		int64
	TeamId		int64
	UserId		int64
	External	bool
	Result		[]*TeamMemberDTO
}
type TeamMemberDTO struct {
	OrgId		int64		`json:"orgId"`
	TeamId		int64		`json:"teamId"`
	UserId		int64		`json:"userId"`
	External	bool		`json:"-"`
	Email		string		`json:"email"`
	Login		string		`json:"login"`
	AvatarUrl	string		`json:"avatarUrl"`
	Labels		[]string	`json:"labels"`
}
