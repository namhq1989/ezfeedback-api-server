package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type RemoveProjectCampaignCategoryHandler struct {
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository
	cachingRepository                 domain.CachingRepository
	service                           domain.Service
}

func NewRemoveProjectCampaignCategoryHandler(
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) RemoveProjectCampaignCategoryHandler {
	return RemoveProjectCampaignCategoryHandler{
		projectCampaignCategoryRepository: projectCampaignCategoryRepository,
		cachingRepository:                 cachingRepository,
		service:                           service,
	}
}

// RemoveProjectCampaignCategory godoc
// @tags     Project
// @summary  Remove project campaign category
// @id       project-remove-campaign-category
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    campaignId  	 path     string true "Campaign id"
// @param    payload body    dto.RemoveProjectCampaignCategoryRequest true "Body"
// @success  200     {object} dto.RemoveProjectCampaignCategoryResponse
// @router   /api/project/{projectId}/campaign/{campaignId}/category [delete]
func (h RemoveProjectCampaignCategoryHandler) RemoveProjectCampaignCategory(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.RemoveProjectCampaignCategoryRequest) (*dto.RemoveProjectCampaignCategoryResponse, error) {
	ctx.Logger().Info("new remove project campaign category request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "campaignID": campaignID, "categoryID": req.CategoryID,
	})

	ctx.Logger().Text("find project campaign in db")
	campaign, err := h.service.GetProjectCampaign(ctx, projectID, campaignID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find project campaign in db", err, appcontext.Fields{})
		return nil, err
	}
	if campaign == nil {
		ctx.Logger().ErrorText("project campaign not found")
		return nil, apperrors.Project.InvalidCampaign
	}

	ctx.Logger().Text("find project category in db")
	category, err := h.service.GetProjectCategory(ctx, projectID, req.CategoryID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find project category in db", err, appcontext.Fields{})
		return nil, err
	}
	if category == nil {
		ctx.Logger().ErrorText("project category not found")
		return nil, apperrors.Project.InvalidCategory
	}

	ctx.Logger().Text("find campaign category in db")
	campaignCategory, err := h.projectCampaignCategoryRepository.FindByCampaignIDAndCategoryID(ctx, campaign.ID, category.ID)
	if err != nil {
		ctx.Logger().Error("failed to find campaign category in db", err, appcontext.Fields{})
		return nil, err
	}
	if campaignCategory == nil {
		ctx.Logger().Text("campaign category does not exist, respond")
		return &dto.RemoveProjectCampaignCategoryResponse{}, nil
	}

	ctx.Logger().Text("delete project campaign category in db")
	if err = h.projectCampaignCategoryRepository.Delete(ctx, *campaignCategory); err != nil {
		ctx.Logger().Error("failed to delete project campaign category in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCampaignsByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done add project campaign category request")
	return &dto.RemoveProjectCampaignCategoryResponse{}, nil
}
