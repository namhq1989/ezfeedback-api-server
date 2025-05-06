package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CountFeedbacksHandler struct {
	feedbackRepository domain.FeedbackRepository
}

func NewCountFeedbacksHandler(feedbackRepository domain.FeedbackRepository) CountFeedbacksHandler {
	return CountFeedbacksHandler{
		feedbackRepository: feedbackRepository,
	}
}

// CountFeedbacks godoc
// @tags     Feedback
// @summary  Count feedbacks
// @id       feedback-count
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.CountFeedbacksRequest true "Query"
// @success  200     {object} dto.CountFeedbacksResponse
// @router   /api/feedback/count [get]
func (h CountFeedbacksHandler) CountFeedbacks(ctx *appcontext.AppContext, performerID string, req dto.CountFeedbacksRequest) (*dto.CountFeedbacksResponse, error) {
	ctx.Logger().Info("new count feedbacks request", appcontext.Fields{
		"performerID": performerID, "projectID": req.ProjectID, "categoryID": req.CategoryID, "campaignType": req.CampaignType,
		"keyword": req.Keyword, "state": req.State, "rating": req.Rating,
	})

	ctx.Logger().Text("create filter")
	filter, err := domain.NewFeedbackFilter(req.ProjectID, req.CampaignType, req.CategoryID, req.Keyword, req.State, req.Rating, 0)
	if err != nil {
		ctx.Logger().Error("failed to create filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count feedbacks in db")
	total, err := h.feedbackRepository.CountWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to count feedbacks in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done count feedbacks request")
	return &dto.CountFeedbacksResponse{Total: total}, nil
}
