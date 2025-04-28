package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

type ProjectSetting struct {
	Domain       string `json:"domain"`
	PrimaryColor string `json:"primaryColor"`
}

func (ProjectSetting) FromDomain(setting domain.ProjectSetting) ProjectSetting {
	return ProjectSetting{
		Domain:       setting.Domain,
		PrimaryColor: setting.PrimaryColor,
	}
}
