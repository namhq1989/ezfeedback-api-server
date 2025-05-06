package query

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetFeedbacksHandler struct {
	feedbackRepository domain.FeedbackRepository
	projectHub         domain.ProjectHub
}

func NewGetFeedbacksHandler(feedbackRepository domain.FeedbackRepository, projectHub domain.ProjectHub) GetFeedbacksHandler {
	return GetFeedbacksHandler{
		feedbackRepository: feedbackRepository,
		projectHub:         projectHub,
	}
}

// GetFeedbacks godoc
// @tags     Feedback
// @summary  Get feedbacks
// @id       feedback-list-all
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.GetFeedbacksRequest true "Query"
// @success  200     {object} dto.GetFeedbacksResponse
// @router   /api/feedback [get]
func (h GetFeedbacksHandler) GetFeedbacks(ctx *appcontext.AppContext, performerID string, req dto.GetFeedbacksRequest) (*dto.GetFeedbacksResponse, error) {
	ctx.Logger().Info("new get feedbacks request", appcontext.Fields{
		"performerID": performerID, "projectID": req.ProjectID, "categoryID": req.CategoryID, "campaignType": req.CampaignType,
		"keyword": req.Keyword, "state": req.State, "rating": req.Rating, "page": req.Page,
	})

	ctx.Logger().Text("get project data via grpc")
	project, err := h.projectHub.GetProjectByID(ctx, req.ProjectID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to get project data via grpc", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().Text("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("create filter")
	filter, err := domain.NewFeedbackFilter(project.ID, req.CampaignType, req.CategoryID, req.Keyword, req.State, req.Rating, req.Page)
	if err != nil {
		ctx.Logger().Error("failed to create filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find feedbacks in db")
	feedbacks, err := h.feedbackRepository.FindWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to find feedbacks in db", err, appcontext.Fields{})
		return nil, err
	}

	if len(feedbacks) == 0 {
		ctx.Logger().Text("no feedbacks found")
		return &dto.GetFeedbacksResponse{
			Feedbacks: make([]dto.Feedback, 0),
		}, nil
	}

	ctx.Logger().Text("convert to dto")
	result := h.convertToDto(ctx, project.ID, feedbacks, filter.Limit)

	ctx.Logger().Text("done get feedbacks request")
	return result, nil
}

func (h GetFeedbacksHandler) convertToDto(ctx *appcontext.AppContext, projectID string, feedbacks []domain.Feedback, limit int64) *dto.GetFeedbacksResponse {
	var result = &dto.GetFeedbacksResponse{
		Feedbacks: make([]dto.Feedback, 0),
		Limit:     limit,
	}

	ctx.Logger().Text("get project campaigns via grpc")
	projectCampaigns, err := h.projectHub.GetProjectCampaigns(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to get project campaigns via grpc", err, appcontext.Fields{})
		return result
	}

	ctx.Logger().Text("get project categories via grpc")
	projectCategories, err := h.projectHub.GetProjectCategories(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to get project categories via grpc", err, appcontext.Fields{})
		return result
	}

	for _, feedback := range feedbacks {
		var (
			campaign = h.findCampaign(projectCampaigns, feedback.CampaignID)
			category = h.findCategory(projectCategories, feedback.CategoryID)
			doc      = dto.Feedback{}.FromDomain(feedback, campaign, category)
		)

		result.Feedbacks = append(result.Feedbacks, doc)
	}

	return result
}

func (h GetFeedbacksHandler) findCampaign(campaigns []domain.ProjectCampaign, campaignID string) domain.ProjectCampaign {
	for _, campaign := range campaigns {
		if campaign.ID == campaignID {
			return campaign
		}
	}
	return domain.ProjectCampaign{
		ID:   "",
		Name: "-",
	}
}

func (h GetFeedbacksHandler) findCategory(categories []domain.ProjectCategory, categoryID string) domain.ProjectCategory {
	for _, category := range categories {
		if category.ID == categoryID {
			return category
		}
	}
	return domain.ProjectCategory{
		ID:   "",
		Name: "-",
	}
}
