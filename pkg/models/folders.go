package models

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrFolderNotFound					= errors.New("Folder not found")
	ErrFolderVersionMismatch			= errors.New("The folder has been changed by someone else")
	ErrFolderTitleEmpty					= errors.New("Folder title cannot be empty")
	ErrFolderWithSameUIDExists			= errors.New("A folder/dashboard with the same uid already exists")
	ErrFolderSameNameExists				= errors.New("A folder or dashboard in the general folder with the same name already exists")
	ErrFolderFailedGenerateUniqueUid	= errors.New("Failed to generate unique folder id")
	ErrFolderAccessDenied				= errors.New("Access denied to folder")
)

type Folder struct {
	Id			int64
	Uid			string
	Title		string
	Url			string
	Version		int
	Created		time.Time
	Updated		time.Time
	UpdatedBy	int64
	CreatedBy	int64
	HasAcl		bool
}

func (cmd *CreateFolderCommand) GetDashboardModel(orgId int64, userId int64) *Dashboard {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dashFolder := NewDashboardFolder(strings.TrimSpace(cmd.Title))
	dashFolder.OrgId = orgId
	dashFolder.SetUid(strings.TrimSpace(cmd.Uid))
	if userId == 0 {
		userId = -1
	}
	dashFolder.CreatedBy = userId
	dashFolder.UpdatedBy = userId
	dashFolder.UpdateSlug()
	return dashFolder
}
func (cmd *UpdateFolderCommand) UpdateDashboardModel(dashFolder *Dashboard, orgId int64, userId int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dashFolder.OrgId = orgId
	dashFolder.Title = strings.TrimSpace(cmd.Title)
	dashFolder.Data.Set("title", dashFolder.Title)
	if cmd.Uid != "" {
		dashFolder.SetUid(cmd.Uid)
	}
	dashFolder.SetVersion(cmd.Version)
	dashFolder.IsFolder = true
	if userId == 0 {
		userId = -1
	}
	dashFolder.UpdatedBy = userId
	dashFolder.UpdateSlug()
}

type CreateFolderCommand struct {
	Uid		string	`json:"uid"`
	Title	string	`json:"title"`
	Result	*Folder
}
type UpdateFolderCommand struct {
	Uid			string	`json:"uid"`
	Title		string	`json:"title"`
	Version		int		`json:"version"`
	Overwrite	bool	`json:"overwrite"`
	Result		*Folder
}
type HasEditPermissionInFoldersQuery struct {
	SignedInUser	*SignedInUser
	Result			bool
}
