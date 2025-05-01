package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/pagetoken"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	feedbackQueryLimit int64 = 20
)

type FeedbackFilter struct {
	ProjectID  string
	CampaignID string
	CategoryID string
	Keyword    string
	State      FeedbackState
	Rating     int32
	Page       int64
	Limit      int64
}

func NewFeedbackFilter(pageToken, projectID, campaignID, categoryID, keyword, state string, rating int32) (*FeedbackFilter, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	if !uuid.IsValidID(campaignID) {
		return nil, apperrors.Project.InvalidCampaign
	}

	if !uuid.IsValidID(categoryID) {
		return nil, apperrors.Project.InvalidCategory
	}

	pt := pagetoken.Decode(pageToken)

	return &FeedbackFilter{
		ProjectID:  projectID,
		CampaignID: campaignID,
		CategoryID: categoryID,
		Keyword:    keyword,
		State:      ToFeedbackState(state),
		Rating:     rating,
		Page:       pt.Page,
		Limit:      feedbackQueryLimit,
	}, nil
}

func IsEndOfFeedbackResult(total int64) bool {
	return total < feedbackQueryLimit
}
