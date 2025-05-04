package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type ProjectHub interface {
	GetProjectByID(ctx *appcontext.AppContext, projectID string) (*Project, error)
	GetProjectCollaborators(ctx *appcontext.AppContext, projectID string) ([]ProjectCollaborator, error)
}

type Project struct {
	ID    string
	Title string
}

type ProjectCollaborator struct {
	ID     string
	UserID string
	Role   ProjectRole
}

type ProjectRole string

const (
	ProjectRoleUnknown ProjectRole = ""
	ProjectRoleOwner   ProjectRole = "owner"
	ProjectRoleEditor  ProjectRole = "editor"
	ProjectRoleViewer  ProjectRole = "viewer"
)

func (s ProjectRole) String() string {
	return string(s)
}

func (s ProjectRole) IsValid() bool {
	return s != ProjectRoleUnknown
}

func (s ProjectRole) IsOwner() bool {
	return s == ProjectRoleOwner
}

func (s ProjectRole) IsEditor() bool {
	return s == ProjectRoleEditor
}

func (s ProjectRole) IsViewer() bool {
	return s == ProjectRoleViewer
}

func ToProjectRole(s string) ProjectRole {
	switch s {
	case ProjectRoleOwner.String():
		return ProjectRoleOwner
	case ProjectRoleEditor.String():
		return ProjectRoleEditor
	case ProjectRoleViewer.String():
		return ProjectRoleViewer
	default:
		return ProjectRoleUnknown
	}
}

type ProjectCampaignType string

const (
	ProjectCampaignTypeUnknown  ProjectCampaignType = ""
	ProjectCampaignTypeFeedback ProjectCampaignType = "feedback"
	ProjectCampaignTypeNPS      ProjectCampaignType = "nps"
	ProjectCampaignTypeCSAT     ProjectCampaignType = "csat"
)

func (s ProjectCampaignType) String() string {
	return string(s)
}

func (s ProjectCampaignType) IsValid() bool {
	return s != ProjectCampaignTypeUnknown
}

func (s ProjectCampaignType) IsFeedback() bool {
	return s == ProjectCampaignTypeFeedback
}

func (s ProjectCampaignType) IsNPS() bool {
	return s == ProjectCampaignTypeNPS
}

func (s ProjectCampaignType) IsCSAT() bool {
	return s == ProjectCampaignTypeCSAT
}

func ToProjectCampaignType(s string) ProjectCampaignType {
	switch s {
	case ProjectCampaignTypeFeedback.String():
		return ProjectCampaignTypeFeedback
	case ProjectCampaignTypeNPS.String():
		return ProjectCampaignTypeNPS
	case ProjectCampaignTypeCSAT.String():
		return ProjectCampaignTypeCSAT
	default:
		return ProjectCampaignTypeUnknown
	}
}
