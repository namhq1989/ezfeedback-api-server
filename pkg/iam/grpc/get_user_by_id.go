package grpc

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetUserByIDHandler struct {
	service domain.Service
}

func NewGetUserByIDHandler(service domain.Service) GetUserByIDHandler {
	return GetUserByIDHandler{
		service: service,
	}
}

func (h GetUserByIDHandler) GetUserById(ctx *appcontext.AppContext, req *iampb.GetUserByIdRequest) (*iampb.GetUserByIdResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get user by id request", appcontext.Fields{"id": req.GetUserId()})

	ctx.Logger().Text("find user in db")
	user, err := h.service.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		ctx.Logger().Error("failed to find user in db", err, appcontext.Fields{})
		return nil, err
	}
	if user == nil {
		ctx.Logger().ErrorText("user not found")
		return nil, apperrors.User.UserNotFound
	}

	var result = &iampb.User{
		Id:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status.String(),
	}

	ctx.Logger().Text("done get user by id request")
	return &iampb.GetUserByIdResponse{
		User: result,
	}, nil
}
