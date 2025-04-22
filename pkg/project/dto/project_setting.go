package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

type ProjectSetting struct {
	IsFeedbackPublic       bool `json:"isFeedbackPublic"`
	AllowAnonymousFeedback bool `json:"allowAnonymousFeedback"`
	EnableVoting           bool `json:"enableVoting"`
}

func (ProjectSetting) FromDomain(setting domain.ProjectSetting) ProjectSetting {
	return ProjectSetting{
		IsFeedbackPublic:       setting.IsFeedbackPublic,
		AllowAnonymousFeedback: setting.AllowAnonymousFeedback,
		EnableVoting:           setting.EnableVoting,
	}
}
