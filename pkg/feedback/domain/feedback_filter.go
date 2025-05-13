package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	feedbackQueryLimit int64 = 10
)

type FeedbackFilter struct {
	ProjectID    string
	CampaignType ProjectCampaignType
	CategoryID   string
	Keyword      string
	State        FeedbackState
	Rating       int32
	Page         int64
	Limit        int64
}

func NewFeedbackFilter(projectID, campaignType, categoryID, keyword, state string, rating int32, page int64) (*FeedbackFilter, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	dCampaignType := ToProjectCampaignType(campaignType)
	if !dCampaignType.IsValid() {
		dCampaignType = ProjectCampaignTypeUnknown
	}

	if categoryID != "" && !uuid.IsValidID(categoryID) {
		return nil, apperrors.Project.InvalidCategory
	}

	return &FeedbackFilter{
		ProjectID:    projectID,
		CampaignType: dCampaignType,
		CategoryID:   categoryID,
		Keyword:      keyword,
		State:        ToFeedbackState(state),
		Rating:       rating,
		Page:         page,
		Limit:        feedbackQueryLimit,
	}, nil
}
