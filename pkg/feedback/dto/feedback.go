package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type Feedback struct {
	ID           string                    `json:"id"`
	Campaign     ProjectCampaign           `json:"campaign"`
	Category     ProjectCategory           `json:"category"`
	AppUserID    string                    `json:"appUserId"`
	Email        string                    `json:"email"`
	Content      string                    `json:"content"`
	Rating       int32                     `json:"rating"`
	IsAnonymous  bool                      `json:"isAnonymous"`
	State        string                    `json:"state"`
	CampaignType string                    `json:"campaignType"`
	Ip           string                    `json:"ip"`
	CountryCode  string                    `json:"countryCode"`
	CreatedAt    *httprespond.TimeResponse `json:"createdAt"`
}

func (Feedback) FromDomain(feedback domain.Feedback, campaign domain.ProjectCampaign, category domain.ProjectCategory) Feedback {
	var appUserID = ""
	if feedback.AppUserID != nil {
		appUserID = *feedback.AppUserID
	}

	var email = ""
	if feedback.Email != nil {
		email = *feedback.Email
	}

	return Feedback{
		ID:           feedback.ID,
		Campaign:     ProjectCampaign{}.FromDomain(campaign),
		Category:     ProjectCategory{}.FromDomain(category),
		AppUserID:    appUserID,
		Email:        email,
		Content:      feedback.Content,
		Rating:       feedback.Rating,
		IsAnonymous:  feedback.IsAnonymous,
		State:        feedback.State.String(),
		CampaignType: feedback.CampaignType.String(),
		Ip:           feedback.Ip,
		CountryCode:  feedback.CountryCode,
		CreatedAt:    httprespond.NewTimeResponse(feedback.CreatedAt),
	}
}
