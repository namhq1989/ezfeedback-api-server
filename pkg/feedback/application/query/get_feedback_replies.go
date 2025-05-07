package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetFeedbackRepliesHandler struct {
	feedbackReplyRepository domain.FeedbackReplyRepository
	iamHub                  domain.IAMHub
}

func NewGetFeedbackRepliesHandler(
	feedbackReplyRepository domain.FeedbackReplyRepository,
	iamHub domain.IAMHub,
) GetFeedbackRepliesHandler {
	return GetFeedbackRepliesHandler{
		feedbackReplyRepository: feedbackReplyRepository,
		iamHub:                  iamHub,
	}
}

// GetFeedbackReplies godoc
// @tags     Feedback
// @summary  Get feedback replies
// @id       feedback-get-replies
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Feedback id"
// @param    payload query    dto.GetFeedbackRepliesRequest true "Query"
// @success  200     {object} dto.GetFeedbackRepliesResponse
// @router   /api/feedback/{id}/reply [get]
func (h GetFeedbackRepliesHandler) GetFeedbackReplies(ctx *appcontext.AppContext, performerID, feedbackID string, req dto.GetFeedbackRepliesRequest) (*dto.GetFeedbackRepliesResponse, error) {
	ctx.Logger().Info("new get feedback replies request", appcontext.Fields{
		"performerID": performerID, "feedbackID": feedbackID, "page": req.Page,
	})

	ctx.Logger().Text("create filter")
	filter, err := domain.NewFeedbackReplyFilter(feedbackID, req.Page)
	if err != nil {
		ctx.Logger().Error("failed to create filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find replies in db")
	replies, err := h.feedbackReplyRepository.FindWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to find replies in db", err, appcontext.Fields{})
		return nil, err
	}
	if len(replies) == 0 {
		ctx.Logger().Text("no replies found, respond")
		return &dto.GetFeedbackRepliesResponse{
			Replies: make([]dto.FeedbackReply, 0),
			Limit:   filter.Limit,
		}, nil
	}

	ctx.Logger().Text("convert to dto")
	result := h.convertToDto(ctx, replies, filter.Limit)

	ctx.Logger().Text("done get feedback replies request")
	return result, nil
}

func (h GetFeedbackRepliesHandler) convertToDto(ctx *appcontext.AppContext, replies []domain.FeedbackReply, limit int64) *dto.GetFeedbackRepliesResponse {
	var result = &dto.GetFeedbackRepliesResponse{
		Replies: make([]dto.FeedbackReply, 0),
		Limit:   limit,
	}

	for _, reply := range replies {
		ctx.Logger().Info("find user via grpc", appcontext.Fields{"userID": reply.UserID})
		user, err := h.iamHub.GetUserByID(ctx, reply.UserID)
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

		doc := dto.FeedbackReply{}.FromDomain(reply, *user)
		result.Replies = append(result.Replies, doc)
	}

	return result
}
