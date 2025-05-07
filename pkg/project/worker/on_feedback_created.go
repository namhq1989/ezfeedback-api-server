package worker

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type OnFeedbackCreatedHandler struct {
	projectRepository         domain.ProjectRepository
	projectCampaignRepository domain.ProjectCampaignRepository
	cachingRepository         domain.CachingRepository
}

func NewOnFeedbackCreatedHandler(
	projectRepository domain.ProjectRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	cachingRepository domain.CachingRepository,
) OnFeedbackCreatedHandler {
	return OnFeedbackCreatedHandler{
		projectRepository:         projectRepository,
		projectCampaignRepository: projectCampaignRepository,
		cachingRepository:         cachingRepository,
	}
}

func (h OnFeedbackCreatedHandler) OnFeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueOnFeedbackCreatedPayload) error {
	if err := h.updateProject(ctx, payload.ProjectID); err != nil {
		return err
	} else {
		return h.updateCampaign(ctx, payload.CampaignID)
	}
}

func (h OnFeedbackCreatedHandler) updateProject(ctx *appcontext.AppContext, projectID string) error {
	ctx.Logger().Text("find project in db")
	project, err := h.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("increase project stats total feedbacks")
	project.AdjustStatsTotalFeedbacks(1)

	ctx.Logger().Text("update project in db")
	if err = h.projectRepository.Update(ctx, *project); err != nil {
		ctx.Logger().Error("failed to update project in db", err, appcontext.Fields{})
		return err
	} else {
		ctx.Logger().Text("set project in caching")
		if err = h.cachingRepository.SetProjectByID(ctx, project.ID, *project); err != nil {
			ctx.Logger().Error("failed to set project in caching", err, appcontext.Fields{})
		}
		if err = h.cachingRepository.DeleteApiGetProjectByID(ctx, project.ID); err != nil {
			ctx.Logger().Error("failed to delete api get project by id in caching data", err, appcontext.Fields{})
		}
	}

	return nil
}

func (h OnFeedbackCreatedHandler) updateCampaign(ctx *appcontext.AppContext, campaignID string) error {
	ctx.Logger().Text("find campaign in db")
	campaign, err := h.projectCampaignRepository.FindByID(ctx, campaignID)
	if err != nil {
		ctx.Logger().Error("failed to find campaign in db", err, appcontext.Fields{})
		return err
	}
	if campaign == nil {
		ctx.Logger().ErrorText("campaign not found")
		return nil
	}

	ctx.Logger().Text("increase campaign stats total feedbacks")
	campaign.AdjustStatsTotalFeedbacks(1)

	ctx.Logger().Text("update campaign in db")
	err = h.projectCampaignRepository.Update(ctx, *campaign)
	return err
}
