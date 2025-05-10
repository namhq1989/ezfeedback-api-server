package query

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectCollaboratorsHandler struct {
	iamHub  domain.IAMHub
	service domain.Service
}

func NewGetProjectCollaboratorsHandler(iamHub domain.IAMHub, service domain.Service) GetProjectCollaboratorsHandler {
	return GetProjectCollaboratorsHandler{
		iamHub:  iamHub,
		service: service,
	}
}

// GetProjectCollaborators godoc
// @tags     Project
// @summary  Get project collaborators
// @id       project-get-collaborators
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload query    dto.GetProjectCollaboratorsRequest true "Query"
// @success  200     {object} dto.GetProjectCollaboratorsResponse
// @router   /api/project/{id}/collaborators [get]
func (h GetProjectCollaboratorsHandler) GetProjectCollaborators(ctx *appcontext.AppContext, performerID, projectID string, _ dto.GetProjectCollaboratorsRequest) (*dto.GetProjectCollaboratorsResponse, error) {
	ctx.Logger().Info("new get project collaborators request", appcontext.Fields{"performerID": performerID, "projectID": projectID})

	ctx.Logger().Text("find project in db")
	project, err := h.service.GetProjectByID(ctx, projectID, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("find project collaborators in db")
	collaborators, err := h.service.GetProjectCollaboratorsByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project collaborators in db", err, appcontext.Fields{})
		return nil, err
	}
	if len(collaborators) == 0 {
		ctx.Logger().Text("this project has no collaborators, respond")
		return &dto.GetProjectCollaboratorsResponse{
			Collaborators: make([]dto.ProjectCollaborator, 0),
		}, nil
	}

	ctx.Logger().Text("convert to dto")
	result := h.convertToDto(ctx, collaborators)

	ctx.Logger().Text("done get project collaborators request")
	return result, nil
}

func (h GetProjectCollaboratorsHandler) convertToDto(ctx *appcontext.AppContext, collaborators []domain.ProjectCollaborator) *dto.GetProjectCollaboratorsResponse {
	var result = make([]dto.ProjectCollaborator, 0)
	for _, collaborator := range collaborators {
		user, err := h.iamHub.GetUserByID(ctx, collaborator.UserID)
		if err != nil {
			ctx.Logger().Error("failed to find user in db", err, appcontext.Fields{})
			continue
		}
		result = append(result, dto.ProjectCollaborator{}.FromDomain(collaborator, *user))
	}

	return &dto.GetProjectCollaboratorsResponse{
		Collaborators: result,
	}
}
