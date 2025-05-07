package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type OnFeedbackCreatedHandler struct {
	queueRepository domain.QueueRepository
}

func NewOnFeedbackCreatedHandler(queueRepository domain.QueueRepository) OnFeedbackCreatedHandler {
	return OnFeedbackCreatedHandler{
		queueRepository: queueRepository,
	}
}

func (h OnFeedbackCreatedHandler) OnFeedbackCreated(ctx *appcontext.AppContext, req *projectpb.OnFeedbackCreatedRequest) (*projectpb.OnFeedbackCreatedResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new on feedback created request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("add task to queue")
	err := h.queueRepository.OnFeedbackCreated(ctx, domain.QueueOnFeedbackCreatedPayload{
		ProjectID:  req.GetProjectId(),
		CampaignID: req.GetCampaignId(),
	})

	ctx.Logger().Text("done on feedback created request")
	return &projectpb.OnFeedbackCreatedResponse{}, err
}
