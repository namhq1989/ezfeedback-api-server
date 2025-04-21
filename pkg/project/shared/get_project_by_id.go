package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectByID(ctx *appcontext.AppContext, id string) (*domain.Project, error) {
	ctx.Logger().Info("[service] get project by id", appcontext.Fields{"id": id})

	ctx.Logger().Text("find project in caching")
	project, err := s.cachingRepository.GetProjectByID(ctx, id)
	if project != nil {
		ctx.Logger().Text("project found in caching, return")
		return project, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project not found in caching, find in db")
	project, err = s.projectRepository.FindByID(ctx, id)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("set project in caching")
	if err = s.cachingRepository.SetProjectByID(ctx, id, *project); err != nil {
		ctx.Logger().Error("failed to set project in caching", err, appcontext.Fields{})
	}
	return project, nil
}
