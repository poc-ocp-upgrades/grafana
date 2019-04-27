package models

import (
	"errors"
)

var (
	ErrPlaylistNotFound		= errors.New("Playlist not found")
	ErrPlaylistWithSameNameExists	= errors.New("A playlist with the same name already exists")
)

type Playlist struct {
	Id		int64	`json:"id"`
	Name		string	`json:"name"`
	Interval	string	`json:"interval"`
	OrgId		int64	`json:"-"`
}
type PlaylistDTO struct {
	Id		int64			`json:"id"`
	Name		string			`json:"name"`
	Interval	string			`json:"interval"`
	OrgId		int64			`json:"-"`
	Items		[]PlaylistItemDTO	`json:"items"`
}
type PlaylistItemDTO struct {
	Id		int64	`json:"id"`
	PlaylistId	int64	`json:"playlistid"`
	Type		string	`json:"type"`
	Title		string	`json:"title"`
	Value		string	`json:"value"`
	Order		int	`json:"order"`
}
type PlaylistDashboard struct {
	Id	int64	`json:"id"`
	Slug	string	`json:"slug"`
	Title	string	`json:"title"`
}
type PlaylistItem struct {
	Id		int64
	PlaylistId	int64
	Type		string
	Value		string
	Order		int
	Title		string
}

func (this PlaylistDashboard) TableName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "dashboard"
}

type Playlists []*Playlist
type PlaylistDashboards []*PlaylistDashboard
type UpdatePlaylistCommand struct {
	OrgId		int64			`json:"-"`
	Id		int64			`json:"id"`
	Name		string			`json:"name" binding:"Required"`
	Interval	string			`json:"interval"`
	Items		[]PlaylistItemDTO	`json:"items"`
	Result		*PlaylistDTO
}
type CreatePlaylistCommand struct {
	Name		string			`json:"name" binding:"Required"`
	Interval	string			`json:"interval"`
	Items		[]PlaylistItemDTO	`json:"items"`
	OrgId		int64			`json:"-"`
	Result		*Playlist
}
type DeletePlaylistCommand struct {
	Id	int64
	OrgId	int64
}
type GetPlaylistsQuery struct {
	Name	string
	Limit	int
	OrgId	int64
	Result	Playlists
}
type GetPlaylistByIdQuery struct {
	Id	int64
	Result	*Playlist
}
type GetPlaylistItemsByIdQuery struct {
	PlaylistId	int64
	Result		*[]PlaylistItem
}
