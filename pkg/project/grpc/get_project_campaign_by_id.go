package grpc

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectCampaignByIDHandler struct {
	projectCampaignHub domain.ProjectCampaignHub
}

func NewGetProjectCampaignByIDHandler(projectCampaignHub domain.ProjectCampaignHub) GetProjectCampaignByIDHandler {
	return GetProjectCampaignByIDHandler{
		projectCampaignHub: projectCampaignHub,
	}
}

func (h GetProjectCampaignByIDHandler) GetProjectCampaignByID(ctx *appcontext.AppContext, req *projectpb.GetProjectCampaignByIdRequest) (*projectpb.GetProjectCampaignByIdResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project campaign by id request", appcontext.Fields{"campaignID": req.GetCampaignId()})

	ctx.Logger().Text("find campaign in db")
	data, err := h.projectCampaignHub.FindProjectCampaignByID(ctx, req.GetCampaignId())
	if err != nil {
		ctx.Logger().Error("failed to find campaign in db", err, appcontext.Fields{})
		return nil, err
	}
	if data == nil {
		ctx.Logger().ErrorText("project campaign not found")
		return nil, apperrors.Project.InvalidCampaign
	}

	var result = &projectpb.GetProjectCampaignByIdResponse{
		Campaign: &projectpb.ProjectCampaign{
			Id:           data.ProjectCampaign.ID,
			CampaignType: data.ProjectCampaign.CampaignType.String(),
			Status:       data.ProjectCampaign.Status.String(),
		},
		Project: &projectpb.Project{
			Id:     data.Project.ID,
			Title:  data.Project.Title,
			Status: data.Project.Status.String(),
			Setting: &projectpb.ProjectSetting{
				Domain: data.ProjectSetting.Domain,
			},
		},
	}

	ctx.Logger().Text("done get project campaign by id request")
	return result, nil
}
