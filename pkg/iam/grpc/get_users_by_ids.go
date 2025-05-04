package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetUsersByIDsHandler struct {
	service domain.Service
}

func NewGetUsersByIDsHandler(service domain.Service) GetUsersByIDsHandler {
	return GetUsersByIDsHandler{
		service: service,
	}
}

func (h GetUsersByIDsHandler) GetUsersByIds(ctx *appcontext.AppContext, req *iampb.GetUsersByIdsRequest) (*iampb.GetUsersByIdsResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get users by ids request", appcontext.Fields{"ids": req.GetUserIds()})

	var result = make([]*iampb.User, 0)

	ctx.Logger().Text("loop each user id")
	for _, userID := range req.GetUserIds() {
		ctx.Logger().Info("find user in db", appcontext.Fields{"userID": userID})
		user, err := h.service.GetUserByID(ctx, userID)
		if err != nil {
			ctx.Logger().Error("failed to find user in db", err, appcontext.Fields{"userID": userID})
			continue
		} else if user == nil {
			ctx.Logger().Error("user not found", nil, appcontext.Fields{"userID": userID})
			continue
		}

		result = append(result, &iampb.User{
			Id:     user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Status: user.Status.String(),
		})
	}

	ctx.Logger().Text("done get users by ids request")
	return &iampb.GetUsersByIdsResponse{
		Users: result,
	}, nil
}
