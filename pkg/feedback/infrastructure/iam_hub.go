package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type IAMHub struct {
	client iampb.IAMServiceClient
}

func NewIAMHub(client iampb.IAMServiceClient) IAMHub {
	return IAMHub{
		client: client,
	}
}

func (r IAMHub) GetUserByID(ctx *appcontext.AppContext, userID string) (*domain.User, error) {
	resp, err := r.client.GetUserById(ctx.Context(), &iampb.GetUserByIdRequest{
		TraceId: ctx.GetTraceID(),
		UserId:  userID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var user = resp.GetUser()
	return &domain.User{
		ID:    user.GetId(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}, nil
}
