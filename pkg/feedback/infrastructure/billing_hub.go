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

func (r BillingHub) CanAcceptFeedback(ctx *appcontext.AppContext, projectID string, totalCreatedFeedbacks int64) (bool, error) {
	resp, err := r.client.CanAcceptFeedback(ctx.Context(), &billingpb.CanAcceptFeedbackRequest{
		TraceId:               ctx.GetTraceID(),
		ProjectId:             projectID,
		TotalCreatedFeedbacks: totalCreatedFeedbacks,
	})
	if err != nil {
		return false, err
	}

	return resp.GetAllowed(), nil
}
