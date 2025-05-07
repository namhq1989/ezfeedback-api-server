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

func (r ProjectHub) GetProjectByID(ctx *appcontext.AppContext, projectID, userID string) (*domain.Project, error) {
	resp, err := r.client.GetProjectById(ctx.Context(), &projectpb.GetProjectByIdRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
		UserId:    userID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetProject() == nil {
		return nil, apperrors.Project.ProjectNotFound
	}

	var project = resp.GetProject()
	return &domain.Project{
		ID:     project.GetId(),
		UserID: project.GetUserId(),
		Title:  project.GetTitle(),
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
			Name:         campaign.GetName(),
			CampaignType: domain.ToProjectCampaignType(campaign.GetCampaignType()),
			Status:       domain.ToStatus(campaign.GetStatus()),
		},
	}, nil
}

func (r ProjectHub) GetProjectCategories(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCategory, error) {
	resp, err := r.client.GetProjectCategories(ctx.Context(), &projectpb.GetProjectCategoriesRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
	})
	if err != nil {
		return nil, err
	}

	var (
		result     = make([]domain.ProjectCategory, 0)
		categories = resp.GetCategories()
	)
	for _, category := range categories {
		result = append(result, domain.ProjectCategory{
			ID:   category.GetId(),
			Name: category.GetName(),
		})
	}

	return result, nil
}

func (r ProjectHub) GetProjectCampaigns(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCampaign, error) {
	resp, err := r.client.GetProjectCampaigns(ctx.Context(), &projectpb.GetProjectCampaignsRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
	})
	if err != nil {
		return nil, err
	}

	var (
		result    = make([]domain.ProjectCampaign, 0)
		campaigns = resp.GetCampaigns()
	)
	for _, campaign := range campaigns {
		result = append(result, domain.ProjectCampaign{
			ID:           campaign.GetId(),
			Name:         campaign.GetName(),
			CampaignType: domain.ToProjectCampaignType(campaign.GetCampaignType()),
			Status:       domain.ToStatus(campaign.GetStatus()),
		})
	}

	return result, nil
}

func (r ProjectHub) OnFeedbackCreated(ctx *appcontext.AppContext, projectID, campaignID string) error {
	_, err := r.client.OnFeedbackCreated(ctx.Context(), &projectpb.OnFeedbackCreatedRequest{
		TraceId:    ctx.GetTraceID(),
		ProjectId:  projectID,
		CampaignId: campaignID,
	})
	return err
}
