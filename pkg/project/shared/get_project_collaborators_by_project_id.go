package shared

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectCollaboratorsByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCollaborator, error) {
	ctx.Logger().Info("[service] get project collaborators by project id", appcontext.Fields{"projectID": projectID})

	ctx.Logger().Text("find project collaborators in caching")
	collaborators, err := s.cachingRepository.GetProjectCollaboratorsByProjectID(ctx, projectID)
	if collaborators != nil {
		ctx.Logger().Text("project collaborators found in caching, return")
		return collaborators, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project collaborators in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project collaborators not found in caching, find in db")
	collaborators, err = s.projectCollaboratorRepository.FindByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project collaborators in db", err, appcontext.Fields{})
		return nil, err
	}
	if collaborators == nil || len(collaborators) == 0 {
		ctx.Logger().Text("project collaborators not found")
		return make([]domain.ProjectCollaborator, 0), nil
	}

	ctx.Logger().Text("set project collaborators in caching")
	if err = s.cachingRepository.SetProjectCollaboratorsByProjectID(ctx, projectID, collaborators); err != nil {
		ctx.Logger().Error("failed to set project collaborators in caching", err, appcontext.Fields{})
	}
	return collaborators, nil
}
