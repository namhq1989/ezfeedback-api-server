package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateProjectCampaignHandler struct {
	projectCampaignRepository domain.ProjectCampaignRepository
	cachingRepository         domain.CachingRepository
	service                   domain.Service
}

func NewUpdateProjectCampaignHandler(
	projectCampaignRepository domain.ProjectCampaignRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) UpdateProjectCampaignHandler {
	return UpdateProjectCampaignHandler{
		projectCampaignRepository: projectCampaignRepository,
		cachingRepository:         cachingRepository,
		service:                   service,
	}
}

// UpdateProjectCampaign godoc
// @tags     Project
// @summary  Update project campaign
// @id       project-update-campaign
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    campaignId  	 path     string true "Campaign id"
// @param    payload body    dto.UpdateProjectCampaignRequest true "Body"
// @success  200     {object} dto.UpdateProjectCampaignResponse
// @router   /api/project/{projectId}/campaign/{campaignId} [put]
func (h UpdateProjectCampaignHandler) UpdateProjectCampaign(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.UpdateProjectCampaignRequest) (*dto.UpdateProjectCampaignResponse, error) {
	ctx.Logger().Info("new update project campaign request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "campaignID": campaignID,
		"name": req.Name, "description": req.Description, "widgetPosition": req.WidgetPosition,
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

	ctx.Logger().Text("set project campaign data")
	if err = campaign.SetName(req.Name); err != nil {
		ctx.Logger().Error("failed to set project campaign name", err, appcontext.Fields{})
		return nil, err
	}
	if err = campaign.SetDescription(req.Description); err != nil {
		ctx.Logger().Error("failed to set project campaign description", err, appcontext.Fields{})
		return nil, err
	}
	if err = campaign.SetSettingWidgetPosition(req.WidgetPosition); err != nil {
		ctx.Logger().Error("failed to set project campaign widget position", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update project campaign in db")
	if err = h.projectCampaignRepository.Update(ctx, *campaign); err != nil {
		ctx.Logger().Error("failed to update project campaign in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCampaignsByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete project campaigns caching data", err, appcontext.Fields{})
	}
	if err = h.cachingRepository.DeleteApiGetProjectByID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete api get project by id caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done update project campaign request")
	return &dto.UpdateProjectCampaignResponse{}, nil
}
