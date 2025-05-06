package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CanAcceptFeedbackHandler struct{}

func NewCanAcceptFeedbackHandler() CanAcceptFeedbackHandler {
	return CanAcceptFeedbackHandler{}
}

func (h CanAcceptFeedbackHandler) CanAcceptFeedback(ctx *appcontext.AppContext, req *billingpb.CanAcceptFeedbackRequest) (*billingpb.CanAcceptFeedbackResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new check can accept request", appcontext.Fields{"projectId": req.ProjectId, "totalCreatedFeedbacks": req.TotalCreatedFeedbacks})

	ctx.Logger().Text("done check can accept request")
	return &billingpb.CanAcceptFeedbackResponse{
		Allowed: true,
	}, nil
}
