package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectCampaign(ctx *appcontext.AppContext, projectID, campaignID, userID string) (*domain.ProjectCampaign, error) {
	ctx.Logger().Info("[service] get project campaign", appcontext.Fields{"projectID": projectID, "campaignID": campaignID, "userID": userID})

	ctx.Logger().Text("find project in db")
	project, err := s.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}
	if !project.IsOwner(userID) {
		ctx.Logger().ErrorText("user is not project owner")
		return nil, apperrors.Common.NotFound
	}

	ctx.Logger().Text("find project campaign in db")
	campaign, err := s.projectCampaignRepository.FindByID(ctx, campaignID)
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

	return campaign, nil
}
