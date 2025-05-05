package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
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

func (h GetProjectCollaboratorsHandler) GetProjectCollaborators(ctx *appcontext.AppContext, req *projectpb.GetProjectCollaboratorsRequest) (*projectpb.GetProjectCollaboratorsResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project collaborators request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("find collaborators in db")
	collaborators, err := h.service.GetProjectCollaboratorsByProjectID(ctx, req.GetProjectId())
	if err != nil {
		ctx.Logger().Error("failed to find collaborators in db", err, appcontext.Fields{})
		return nil, err
	}

	var result = make([]*projectpb.ProjectCollaborator, 0)
	for _, collaborator := range collaborators {
		user, uErr := h.iamHub.GetUserByID(ctx, collaborator.UserID)
		if uErr != nil {
			ctx.Logger().Error("failed to find user in db", uErr, appcontext.Fields{})
			continue
		}

		result = append(result, &projectpb.ProjectCollaborator{
			Id: collaborator.ID,
			User: &projectpb.Collaborator{
				Id:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
			Role: collaborator.Role.String(),
		})
	}

	ctx.Logger().Text("done get project collaborators request")
	return &projectpb.GetProjectCollaboratorsResponse{
		Collaborators: result,
	}, nil
}
