package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetFeedbackStateHistoriesHandler struct {
	feedbackStateHistoryRepository domain.FeedbackStateHistoryRepository
	iamHub                         domain.IAMHub
}

func NewGetFeedbackStateHistoriesHandler(
	feedbackStateHistoryRepository domain.FeedbackStateHistoryRepository,
	iamHub domain.IAMHub,
) GetFeedbackStateHistoriesHandler {
	return GetFeedbackStateHistoriesHandler{
		feedbackStateHistoryRepository: feedbackStateHistoryRepository,
		iamHub:                         iamHub,
	}
}

// GetFeedbackStateHistories godoc
// @tags     Feedback
// @summary  Get feedback state histories
// @id       feedback-get-state-histories
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Feedback id"
// @param    payload query    dto.GetFeedbackStateHistoriesRequest true "Query"
// @success  200     {object} dto.GetFeedbackStateHistoriesResponse
// @router   /api/feedback/{id}/state-history [get]
func (h GetFeedbackStateHistoriesHandler) GetFeedbackStateHistories(ctx *appcontext.AppContext, performerID, feedbackID string, req dto.GetFeedbackStateHistoriesRequest) (*dto.GetFeedbackStateHistoriesResponse, error) {
	ctx.Logger().Info("new get feedback state histories request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "page": req.Page,
	})

	ctx.Logger().Text("create filter")
	filter, err := domain.NewFeedbackStateHistoryFilter(feedbackID, req.Page)
	if err != nil {
		ctx.Logger().Error("failed to create filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find state histories in db")
	histories, err := h.feedbackStateHistoryRepository.FindWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to find state histories in db", err, appcontext.Fields{})
		return nil, err
	}
	if len(histories) == 0 {
		ctx.Logger().Text("no state histories found, respond")
		return &dto.GetFeedbackStateHistoriesResponse{
			Histories: make([]dto.FeedbackStateHistory, 0),
			Limit:     filter.Limit,
		}, nil
	}

	ctx.Logger().Text("convert to dto")
	result := h.convertToDto(ctx, histories, filter.Limit)

	ctx.Logger().Text("done get feedback state histories request")
	return result, nil
}

func (h GetFeedbackStateHistoriesHandler) convertToDto(ctx *appcontext.AppContext, histories []domain.FeedbackStateHistory, limit int64) *dto.GetFeedbackStateHistoriesResponse {
	var result = &dto.GetFeedbackStateHistoriesResponse{
		Histories: make([]dto.FeedbackStateHistory, 0),
		Limit:     limit,
	}

	for _, history := range histories {
		ctx.Logger().Info("find user via grpc", appcontext.Fields{"userID": history.UserID})
		user, err := h.iamHub.GetUserByID(ctx, history.UserID)
		if err != nil {
			ctx.Logger().Error("failed to find user via grpc", err, appcontext.Fields{})
			return result
		}
		if user == nil {
			ctx.Logger().ErrorText("user not found, create an empty one")
			user = &domain.User{
				ID:    "",
				Name:  "-",
				Email: "-",
			}
		}

		doc := dto.FeedbackStateHistory{}.FromDomain(history, *user)
		result.Histories = append(result.Histories, doc)
	}

	return result
}
