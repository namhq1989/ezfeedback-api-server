package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type AddProjectCampaignCategoryHandler struct {
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository
	cachingRepository                 domain.CachingRepository
	service                           domain.Service
}

func NewAddProjectCampaignCategoryHandler(
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) AddProjectCampaignCategoryHandler {
	return AddProjectCampaignCategoryHandler{
		projectCampaignCategoryRepository: projectCampaignCategoryRepository,
		cachingRepository:                 cachingRepository,
		service:                           service,
	}
}

// AddProjectCampaignCategory godoc
// @tags     Project
// @summary  Add project campaign category
// @id       project-add-campaign-category
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    campaignId  	 path     string true "Campaign id"
// @param    payload body    dto.AddProjectCampaignCategoryRequest true "Body"
// @success  200     {object} dto.AddProjectCampaignCategoryResponse
// @router   /api/project/{projectId}/campaign/{campaignId}/category [post]
func (h AddProjectCampaignCategoryHandler) AddProjectCampaignCategory(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.AddProjectCampaignCategoryRequest) (*dto.AddProjectCampaignCategoryResponse, error) {
	ctx.Logger().Info("new add project campaign category request", appcontext.Fields{
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
	if campaignCategory != nil {
		ctx.Logger().Text("campaign category already exists, respond")
		return &dto.AddProjectCampaignCategoryResponse{}, nil
	}

	ctx.Logger().Text("new project campaign category model")
	campaignCategory, err = domain.NewProjectCampaignCategory(campaign.ID, category.ID)
	if err != nil {
		ctx.Logger().Error("failed to create project campaign category model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist project campaign category in db")
	if err = h.projectCampaignCategoryRepository.Create(ctx, *campaignCategory); err != nil {
		ctx.Logger().Error("failed to persist project campaign category in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCampaignsByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete project campaigns caching data", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectByID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete api get project by id caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done add project campaign category request")
	return &dto.AddProjectCampaignCategoryResponse{}, nil
}
