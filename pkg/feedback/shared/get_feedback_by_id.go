package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetFeedbackByID(ctx *appcontext.AppContext, feedbackID string) (*domain.Feedback, error) {
	ctx.Logger().Info("[service] get feedback by id", appcontext.Fields{"feedbackID": feedbackID})

	ctx.Logger().Text("find feedback in caching")
	feedback, err := s.cachingRepository.GetFeedbackByID(ctx, feedbackID)
	if feedback != nil {
		ctx.Logger().Text("feedback found in caching, return")
		return feedback, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find feedback in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("feedback not found in caching, find in db")
	feedback, err = s.feedbackRepository.FindByID(ctx, feedbackID)
	if err != nil {
		ctx.Logger().Error("failed to find feedback in db", err, appcontext.Fields{})
		return nil, err
	}
	if feedback == nil {
		ctx.Logger().ErrorText("feedback not found")
		return nil, apperrors.Feedback.FeedbackNotFound
	}

	ctx.Logger().Text("set feedback in caching")
	if err = s.cachingRepository.SetFeedbackByID(ctx, feedbackID, *feedback); err != nil {
		ctx.Logger().Error("failed to set feedback in caching", err, appcontext.Fields{})
	}
	return feedback, nil
}
