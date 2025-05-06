package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectCampaignsHandler struct {
	service domain.Service
}

func NewGetProjectCampaignsHandler(service domain.Service) GetProjectCampaignsHandler {
	return GetProjectCampaignsHandler{
		service: service,
	}
}

func (h GetProjectCampaignsHandler) GetProjectCampaigns(ctx *appcontext.AppContext, req *projectpb.GetProjectCampaignsRequest) (*projectpb.GetProjectCampaignsResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project campaigns request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("find campaigns in db")
	campaigns, err := h.service.GetProjectCampaignsByProjectID(ctx, req.GetProjectId())
	if err != nil {
		ctx.Logger().Error("failed to find campaigns in db", err, appcontext.Fields{})
		return nil, err
	}

	var result = make([]*projectpb.ProjectCampaign, 0)
	for _, campaign := range campaigns {
		result = append(result, &projectpb.ProjectCampaign{
			Id:           campaign.ID,
			Name:         campaign.Name,
			CampaignType: campaign.CampaignType.String(),
			Status:       campaign.Status.String(),
		})
	}

	ctx.Logger().Text("done get project campaigns request")
	return &projectpb.GetProjectCampaignsResponse{
		Campaigns: result,
	}, nil
}
