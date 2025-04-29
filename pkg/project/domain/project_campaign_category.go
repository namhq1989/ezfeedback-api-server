package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type ProjectCampaignCategoryRepository interface {
	Create(ctx *appcontext.AppContext, campaignCategory ProjectCampaignCategory) error
	Delete(ctx *appcontext.AppContext, campaignCategory ProjectCampaignCategory) error
	FindByID(ctx *appcontext.AppContext, campaignCategoryID string) (*ProjectCampaignCategory, error)
	FindByCampaignIDAndCategoryID(ctx *appcontext.AppContext, campaignID string, categoryID string) (*ProjectCampaignCategory, error)
	FindByCampaignID(ctx *appcontext.AppContext, campaignID string) ([]ProjectCampaignCategory, error)
}

type ProjectCampaignCategory struct {
	ID         string
	CampaignID string
	CategoryID string
	CreatedAt  time.Time
}

func NewProjectCampaignCategory(campaignID, categoryID string) (*ProjectCampaignCategory, error) {
	if !uuid.IsValidID(campaignID) {
		return nil, apperrors.Project.InvalidCampaign
	}

	if !uuid.IsValidID(categoryID) {
		return nil, apperrors.Project.InvalidCategory
	}

	var now = manipulation.NowUTC()
	return &ProjectCampaignCategory{
		ID:         uuid.New(),
		CampaignID: campaignID,
		CategoryID: categoryID,
		CreatedAt:  now,
	}, nil
}
