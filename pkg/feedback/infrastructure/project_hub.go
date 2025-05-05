package infrastructure

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ProjectHub struct {
	client projectpb.ProjectServiceClient
}

func NewProjectHub(client projectpb.ProjectServiceClient) ProjectHub {
	return ProjectHub{
		client: client,
	}
}

func (r ProjectHub) GetProjectByID(ctx *appcontext.AppContext, projectID string) (*domain.Project, error) {
	resp, err := r.client.GetProjectById(ctx.Context(), &projectpb.GetProjectByIdRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetProject() == nil {
		return nil, apperrors.Project.ProjectNotFound
	}

	var project = resp.GetProject()
	return &domain.Project{
		ID:    project.GetId(),
		Title: project.GetTitle(),
	}, nil
}

func (r ProjectHub) GetProjectCampaignByID(ctx *appcontext.AppContext, campaignID string) (*domain.ProjectCampaignHubData, error) {
	resp, err := r.client.GetProjectCampaignById(ctx.Context(), &projectpb.GetProjectCampaignByIdRequest{
		TraceId:    ctx.GetTraceID(),
		CampaignId: campaignID,
	})
	if err != nil {
		return nil, err
	}

	campaign := resp.GetCampaign()
	if campaign == nil {
		return nil, apperrors.Project.InvalidCampaign
	}

	project := resp.GetProject()
	if project == nil {
		return nil, apperrors.Project.ProjectNotFound
	}

	return &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     project.GetId(),
			Status: domain.ToStatus(project.GetStatus()),
			Setting: domain.ProjectSetting{
				Domain: project.GetSetting().GetDomain(),
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaign.GetId(),
			CampaignType: domain.ToProjectCampaignType(campaign.GetCampaignType()),
			Status:       domain.ToStatus(campaign.GetStatus()),
		},
	}, nil
}
