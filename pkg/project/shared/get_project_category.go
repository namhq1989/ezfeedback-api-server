package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectCategory(ctx *appcontext.AppContext, projectID, categoryID, userID string) (*domain.ProjectCategory, error) {
	ctx.Logger().Info("[service] get project category", appcontext.Fields{"projectID": projectID, "categoryID": categoryID, "userID": userID})

	ctx.Logger().Text("find project in db")
	project, err := s.GetProjectByID(ctx, projectID, userID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("find project category in db")
	category, err := s.projectCategoryRepository.FindByID(ctx, categoryID)
	if err != nil {
		ctx.Logger().Error("failed to find project category in db", err, appcontext.Fields{})
		return nil, err
	}
	if category == nil {
		ctx.Logger().ErrorText("project category not found")
		return nil, apperrors.Project.InvalidCategory
	}
	if !category.IsBelongToProject(projectID) {
		ctx.Logger().ErrorText("project category not belong to project")
		return nil, apperrors.Project.InvalidCategory
	}

	return category, nil
}
