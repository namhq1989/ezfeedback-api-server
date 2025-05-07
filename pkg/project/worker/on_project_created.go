package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type OnProjectCreatedHandler struct {
	projectCollaboratorRepository domain.ProjectCollaboratorRepository
}

func NewOnProjectCreatedHandler(
	projectCollaboratorRepository domain.ProjectCollaboratorRepository,
) OnProjectCreatedHandler {
	return OnProjectCreatedHandler{
		projectCollaboratorRepository: projectCollaboratorRepository,
	}
}

func (h OnProjectCreatedHandler) OnProjectCreated(ctx *appcontext.AppContext, payload domain.QueueOnProjectCreatedPayload) error {
	ctx.Logger().Text("create project collaborator")
	collaborator, err := domain.NewProjectCollaborator(payload.ProjectID, payload.UserID, domain.ProjectRoleOwner.String())
	if err != nil {
		ctx.Logger().Error("failed to create project collaborator", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("persist project collaborator to db")
	err = h.projectCollaboratorRepository.Create(ctx, *collaborator)
	return err
}
