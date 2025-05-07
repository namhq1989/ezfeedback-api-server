package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateFeedbackReplyHandler struct {
	feedbackReplyRepository domain.FeedbackReplyRepository
	feedbackRepository      domain.FeedbackRepository
	service                 domain.Service
}

func NewCreateFeedbackReplyHandler(
	feedbackReplyRepository domain.FeedbackReplyRepository,
	feedbackRepository domain.FeedbackRepository,
	service domain.Service,
) CreateFeedbackReplyHandler {
	return CreateFeedbackReplyHandler{
		feedbackReplyRepository: feedbackReplyRepository,
		feedbackRepository:      feedbackRepository,
		service:                 service,
	}
}

// CreateFeedbackReply godoc
// @tags     Feedback
// @summary  Create feedback reply
// @id       feedback-create-reply
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Feedback id"
// @param    payload body    dto.CreateFeedbackReplyRequest true "Body"
// @success  200     {object} dto.CreateFeedbackReplyResponse
// @router   /api/feedback/{id}/reply [post]
func (h CreateFeedbackReplyHandler) CreateFeedbackReply(ctx *appcontext.AppContext, performerID, feedbackID string, req dto.CreateFeedbackReplyRequest) (*dto.CreateFeedbackReplyResponse, error) {
	ctx.Logger().Info("new create feedback reply request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "content": req.Content,
	})

	ctx.Logger().Text("find feedback in db")
	feedback, err := h.service.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		ctx.Logger().Error("failed to find feedback in db", err, appcontext.Fields{})
		return nil, err
	}
	if feedback == nil {
		ctx.Logger().ErrorText("feedback not found")
		return nil, apperrors.Feedback.FeedbackNotFound
	}

	ctx.Logger().Text("create reply model")
	reply, err := domain.NewFeedbackReply(feedbackID, performerID, req.Content)
	if err != nil {
		ctx.Logger().Error("failed to create reply model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist reply to db")
	if err = h.feedbackReplyRepository.Create(ctx, *reply); err != nil {
		ctx.Logger().Error("failed to persist reply to db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("increase feedback total reply stats")
	feedback.AdjustStatsTotalReplies(1)

	ctx.Logger().Text("update feedback in db")
	if err = h.feedbackRepository.Update(ctx, *feedback); err != nil {
		ctx.Logger().Error("failed to update feedback in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done create feedback reply")
	return &dto.CreateFeedbackReplyResponse{}, nil
}
