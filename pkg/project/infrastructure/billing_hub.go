package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/go-utilities/appcontext"
)

type BillingHub struct {
	client billingpb.BillingServiceClient
}

func NewBillingHub(client billingpb.BillingServiceClient) BillingHub {
	return BillingHub{
		client: client,
	}
}

func (r BillingHub) CanCreateProject(ctx *appcontext.AppContext, userID string) (bool, error) {
	resp, err := r.client.CanCreateProject(ctx.Context(), &billingpb.CanCreateProjectRequest{
		TraceId: ctx.GetTraceID(),
		UserId:  userID,
	})
	if err != nil {
		return false, err
	}

	return resp.GetAllowed(), nil
}
