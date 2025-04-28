package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ChangeProjectCampaignStatusHandler struct {
	projectCampaignRepository domain.ProjectCampaignRepository
	cachingRepository         domain.CachingRepository
	service                   domain.Service
}

func NewChangeProjectCampaignStatusHandler(
	projectCampaignRepository domain.ProjectCampaignRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) ChangeProjectCampaignStatusHandler {
	return ChangeProjectCampaignStatusHandler{
		projectCampaignRepository: projectCampaignRepository,
		cachingRepository:         cachingRepository,
		service:                   service,
	}
}

// ChangeProjectCampaignStatus godoc
// @tags     Project
// @summary  Change project campaign status
// @id       project-change-campaign-status
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    projectId  	 path     string true "Project id"
// @param    campaignId  	 path     string true "Campaign id"
// @param    payload body    dto.ChangeProjectCampaignStatusRequest true "Body"
// @success  200     {object} dto.ChangeProjectCampaignStatusResponse
// @router   /api/project/{projectId}/campaign/{campaignId}/status [patch]
func (h ChangeProjectCampaignStatusHandler) ChangeProjectCampaignStatus(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.ChangeProjectCampaignStatusRequest) (*dto.ChangeProjectCampaignStatusResponse, error) {
	ctx.Logger().Info("new change project campaign status request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "campaignID": campaignID, "status": req.Status,
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

	if campaign.Status.IsEqual(req.Status) {
		ctx.Logger().Text("project campaign status not changed, respond")
		return &dto.ChangeProjectCampaignStatusResponse{}, nil
	}

	ctx.Logger().Text("set project campaign new status")
	if err = campaign.SetStatus(req.Status); err != nil {
		ctx.Logger().Error("failed to set project campaign new status", err, appcontext.Fields{})
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

	ctx.Logger().Text("done change project category status request")
	return &dto.ChangeProjectCampaignStatusResponse{}, nil
}
