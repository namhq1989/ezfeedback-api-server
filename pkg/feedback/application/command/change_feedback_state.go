package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ChangeFeedbackStateHandler struct {
	feedbackRepository             domain.FeedbackRepository
	feedbackStateHistoryRepository domain.FeedbackStateHistoryRepository
	cachingRepository              domain.CachingRepository
	service                        domain.Service
}

func NewChangeFeedbackStateHandler(
	feedbackRepository domain.FeedbackRepository,
	feedbackStateHistoryRepository domain.FeedbackStateHistoryRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) ChangeFeedbackStateHandler {
	return ChangeFeedbackStateHandler{
		feedbackRepository:             feedbackRepository,
		feedbackStateHistoryRepository: feedbackStateHistoryRepository,
		cachingRepository:              cachingRepository,
		service:                        service,
	}
}

// ChangeFeedbackState godoc
// @tags     Feedback
// @summary  Change feedback state
// @id       feedback-change-state
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Feedback id"
// @param    payload body    dto.ChangeFeedbackStateRequest true "Body"
// @success  200     {object} dto.ChangeFeedbackStateResponse
// @router   /api/feedback/{id}/state [patch]
func (h ChangeFeedbackStateHandler) ChangeFeedbackState(ctx *appcontext.AppContext, performerID, feedbackID string, req dto.ChangeFeedbackStateRequest) (*dto.ChangeFeedbackStateResponse, error) {
	ctx.Logger().Info("new change feedback state request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "state": req.State,
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

	if feedback.State.IsEqual(req.State) {
		ctx.Logger().Text("feedback state not changed, skip")
		return &dto.ChangeFeedbackStateResponse{}, nil
	}

	ctx.Logger().Text("update feedback state")
	if err = feedback.SetState(req.State); err != nil {
		ctx.Logger().Error("failed to update feedback state", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update feedback in db")
	if err = h.feedbackRepository.Update(ctx, *feedback); err != nil {
		ctx.Logger().Error("failed to update feedback in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update feedback data in caching")
	if err = h.cachingRepository.SetFeedbackByID(ctx, feedbackID, *feedback); err != nil {
		ctx.Logger().Error("failed to update feedback data in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("create feedback state history")
	history, err := domain.NewFeedbackStateHistory(feedbackID, performerID, req.State)
	if err != nil {
		ctx.Logger().Error("failed to create feedback state history", err, appcontext.Fields{})
	} else {
		ctx.Logger().Text("persist state history to db")
		if err = h.feedbackStateHistoryRepository.Create(ctx, *history); err != nil {
			ctx.Logger().Error("failed to persist state history to db", err, appcontext.Fields{})
		}
	}

	ctx.Logger().Text("done change feedback state")
	return &dto.ChangeFeedbackStateResponse{}, nil
}
