package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
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

func (r IAMHub) GetUsersByIDs(ctx *appcontext.AppContext, userIDs []string) ([]domain.User, error) {
	var result = make([]domain.User, 0)

	resp, err := r.client.GetUsersByIds(ctx.Context(), &iampb.GetUsersByIdsRequest{
		TraceId: ctx.GetTraceID(),
		UserIds: userIDs,
	})
	if err != nil {
		return result, err
	}

	for _, user := range resp.GetUsers() {
		result = append(result, domain.User{
			ID:     user.GetId(),
			Name:   user.GetName(),
			Email:  user.GetEmail(),
			Status: domain.UserStatus(user.GetStatus()),
		})
	}

	return result, nil
}
