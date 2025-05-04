package infrastructure

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
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

func (r ProjectHub) GetProjectCollaborators(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCollaborator, error) {
	var result = make([]domain.ProjectCollaborator, 0)

	resp, err := r.client.GetProjectCollaborators(ctx.Context(), &projectpb.GetProjectCollaboratorsRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
	})
	if err != nil {
		return result, err
	}

	for _, collaborator := range resp.GetCollaborators() {
		result = append(result, domain.ProjectCollaborator{
			ID:     collaborator.GetId(),
			UserID: collaborator.GetUserId(),
			Role:   domain.ProjectRole(collaborator.GetRole()),
		})
	}

	return result, nil
}
