package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type AddProjectCampaignCategoryHandler struct {
	projectRepository                 domain.ProjectRepository
	projectCampaignRepository         domain.ProjectCampaignRepository
	projectCategoryRepository         domain.ProjectCategoryRepository
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository
	cachingRepository                 domain.CachingRepository
}

func NewAddProjectCampaignCategoryHandler(
	projectRepository domain.ProjectRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	projectCategoryRepository domain.ProjectCategoryRepository,
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository,
	cachingRepository domain.CachingRepository,
) AddProjectCampaignCategoryHandler {
	return AddProjectCampaignCategoryHandler{
		projectRepository:                 projectRepository,
		projectCampaignRepository:         projectCampaignRepository,
		projectCategoryRepository:         projectCategoryRepository,
		projectCampaignCategoryRepository: projectCampaignCategoryRepository,
		cachingRepository:                 cachingRepository,
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

	ctx.Logger().Text("find project in db")
	project, err := h.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}
	if !project.IsOwner(performerID) {
		ctx.Logger().ErrorText("user is not project owner")
		return nil, apperrors.Common.NotFound
	}

	ctx.Logger().Text("find project campaign in db")
	campaign, err := h.projectCampaignRepository.FindByID(ctx, campaignID)
	if err != nil {
		ctx.Logger().Error("failed to find project campaign in db", err, appcontext.Fields{})
		return nil, err
	}
	if campaign == nil {
		ctx.Logger().ErrorText("project campaign not found")
		return nil, apperrors.Project.InvalidCampaign
	}
	if !campaign.IsBelongToProject(projectID) {
		ctx.Logger().ErrorText("project campaign not belong to project")
		return nil, apperrors.Project.InvalidCampaign
	}

	// TODO:
	// - create a service to get project campaign by id (project_id, campaign_id)
	// - reuse here
	// - same for the category
	// - fix tests too

	return nil, nil
}
