package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

type ProjectCampaignSetting struct {
	WidgetPosition   string `json:"widgetPosition"`
	AllowAnonymous   bool   `json:"allowAnonymous"`
	EnableRating     bool   `json:"enableRating"`
	MinRating        int32  `json:"minRating"`
	MaxRating        int32  `json:"maxRating"`
	FollowUpQuestion string `json:"followUpQuestion"`
}

func (ProjectCampaignSetting) FromDomain(setting domain.ProjectCampaignSetting) ProjectCampaignSetting {
	return ProjectCampaignSetting{
		WidgetPosition:   setting.WidgetPosition,
		AllowAnonymous:   setting.AllowAnonymous,
		EnableRating:     setting.EnableRating,
		MinRating:        setting.MinRating,
		MaxRating:        setting.MaxRating,
		FollowUpQuestion: setting.FollowUpQuestion,
	}
}
