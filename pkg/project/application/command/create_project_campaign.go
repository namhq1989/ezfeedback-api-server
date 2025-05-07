package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateProjectCampaignHandler struct {
	projectRepository         domain.ProjectRepository
	projectCampaignRepository domain.ProjectCampaignRepository
	cachingRepository         domain.CachingRepository
}

func NewCreateProjectCampaignHandler(
	projectRepository domain.ProjectRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	cachingRepository domain.CachingRepository,
) CreateProjectCampaignHandler {
	return CreateProjectCampaignHandler{
		projectRepository:         projectRepository,
		projectCampaignRepository: projectCampaignRepository,
		cachingRepository:         cachingRepository,
	}
}

// CreateProjectCampaign godoc
// @tags     Project
// @summary  Create project campaign
// @id       project-create-campaign
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload body    dto.CreateProjectCampaignRequest true "Body"
// @success  200     {object} dto.CreateProjectCampaignResponse
// @router   /api/project/{id}/campaign [post]
func (h CreateProjectCampaignHandler) CreateProjectCampaign(ctx *appcontext.AppContext, performerID, projectID string, req dto.CreateProjectCampaignRequest) (*dto.CreateProjectCampaignResponse, error) {
	ctx.Logger().Info("new create project campaign request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "name": req.Name,
		"campaignType": req.CampaignType, "widgetPosition": req.Settings.WidgetPosition,
		"allowAnonymous": req.Settings.AllowAnonymous, "enableRating": req.Settings.EnableRating, "followUpQuestion": req.Settings.FollowUpQuestion,
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

	ctx.Logger().Text("create project campaign model")
	campaign, err := domain.NewProjectCampaign(projectID, req.Name, req.CampaignType, domain.ProjectCampaignSetting{
		WidgetPosition:   req.Settings.WidgetPosition,
		AllowAnonymous:   req.Settings.AllowAnonymous,
		EnableRating:     req.Settings.EnableRating,
		FollowUpQuestion: req.Settings.FollowUpQuestion,
	})
	if err != nil {
		ctx.Logger().Error("failed to create project campaign model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count by campaign type")
	totalCampaign, err := h.projectCampaignRepository.CountTotalByProjectIDAndCampaignType(ctx, projectID, campaign.CampaignType.String())
	if err != nil {
		ctx.Logger().Error("failed to count by campaign type", err, appcontext.Fields{})
		return nil, err
	}
	if domain.IsReachedMaxCampaignPerType(totalCampaign) {
		ctx.Logger().ErrorText("reached max campaign per type")
		return nil, apperrors.Project.CampaignTypeLimitExceeded
	}

	ctx.Logger().Text("persist project campaign in db")
	if err = h.projectCampaignRepository.Create(ctx, *campaign); err != nil {
		ctx.Logger().Error("failed to persist project campaign in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCampaignsByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete project campaigns caching data", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectByID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete api get project by id caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done create project campaign request")
	return &dto.CreateProjectCampaignResponse{ID: campaign.ID}, nil
}
