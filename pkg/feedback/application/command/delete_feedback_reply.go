package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type DeleteFeedbackReplyHandler struct {
	feedbackReplyRepository domain.FeedbackReplyRepository
	feedbackRepository      domain.FeedbackRepository
	service                 domain.Service
}

func NewDeleteFeedbackReplyHandler(
	feedbackReplyRepository domain.FeedbackReplyRepository,
	feedbackRepository domain.FeedbackRepository,
	service domain.Service,
) DeleteFeedbackReplyHandler {
	return DeleteFeedbackReplyHandler{
		feedbackReplyRepository: feedbackReplyRepository,
		feedbackRepository:      feedbackRepository,
		service:                 service,
	}
}

// DeleteFeedbackReply godoc
// @tags     Feedback
// @summary  Delete feedback reply
// @id       feedback-delete-reply
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    feedbackId  	 path     string true "Feedback id"
// @param    replyId  	 path     string true "Reply id"
// @param    payload body    dto.DeleteFeedbackReplyRequest true "Body"
// @success  200     {object} dto.DeleteFeedbackReplyResponse
// @router   /api/feedback/{feedbackId}/reply/{replyId} [delete]
func (h DeleteFeedbackReplyHandler) DeleteFeedbackReply(ctx *appcontext.AppContext, performerID, feedbackID, replyID string, _ dto.DeleteFeedbackReplyRequest) (*dto.DeleteFeedbackReplyResponse, error) {
	ctx.Logger().Info("new delete feedback reply request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "replyID": replyID,
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

	ctx.Logger().Text("find reply in db")
	reply, err := h.feedbackReplyRepository.FindByID(ctx, replyID)
	if err != nil {
		ctx.Logger().Error("failed to find reply in db", err, appcontext.Fields{})
		return nil, err
	}
	if reply == nil {
		ctx.Logger().ErrorText("reply not found")
		return nil, apperrors.Feedback.InvalidReply
	}
	if !reply.IsBelongToFeedback(feedbackID) {
		ctx.Logger().ErrorText("reply not belong to feedback")
		return nil, apperrors.Feedback.InvalidReply
	}

	ctx.Logger().Text("delete reply in db")
	if err = h.feedbackReplyRepository.Delete(ctx, *reply); err != nil {
		ctx.Logger().Error("failed to delete reply in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("decrease feedback total reply stats")
	feedback.AdjustStatsTotalReplies(-1)

	ctx.Logger().Text("update feedback in db")
	if err = h.feedbackRepository.Update(ctx, *feedback); err != nil {
		ctx.Logger().Error("failed to update feedback in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done delete feedback reply")
	return &dto.DeleteFeedbackReplyResponse{}, nil
}
