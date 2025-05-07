package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateFeedbackReplyHandler struct {
	feedbackReplyRepository domain.FeedbackReplyRepository
}

func NewUpdateFeedbackReplyHandler(
	feedbackReplyRepository domain.FeedbackReplyRepository,
) UpdateFeedbackReplyHandler {
	return UpdateFeedbackReplyHandler{
		feedbackReplyRepository: feedbackReplyRepository,
	}
}

// UpdateFeedbackReply godoc
// @tags     Feedback
// @summary  Update feedback reply
// @id       feedback-update-reply
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    feedbackId  	 path     string true "Feedback id"
// @param    replyId  	 path     string true "Reply id"
// @param    payload body    dto.UpdateFeedbackReplyRequest true "Body"
// @success  200     {object} dto.UpdateFeedbackReplyResponse
// @router   /api/feedback/{feedbackId}/reply/{replyId} [put]
func (h UpdateFeedbackReplyHandler) UpdateFeedbackReply(ctx *appcontext.AppContext, performerID, feedbackID, replyID string, req dto.UpdateFeedbackReplyRequest) (*dto.UpdateFeedbackReplyResponse, error) {
	ctx.Logger().Info("new update feedback reply request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "replyID": replyID, "content": req.Content,
	})

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

	ctx.Logger().Text("update reply data")
	if err = reply.SetContent(req.Content); err != nil {
		ctx.Logger().Error("failed to update reply data", err, appcontext.Fields{})
		return nil, err
	}
	reply.MarkAsEdited()

	ctx.Logger().Text("update reply in db")
	if err = h.feedbackReplyRepository.Update(ctx, *reply); err != nil {
		ctx.Logger().Error("failed to update reply in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done update feedback reply")
	return &dto.UpdateFeedbackReplyResponse{}, nil
}
