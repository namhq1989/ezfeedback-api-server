package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateProjectCategoryHandler struct {
	projectRepository         domain.ProjectRepository
	projectCategoryRepository domain.ProjectCategoryRepository
	cachingRepository         domain.CachingRepository
}

func NewCreateProjectCategoryHandler(
	projectRepository domain.ProjectRepository,
	projectCategoryRepository domain.ProjectCategoryRepository,
	cachingRepository domain.CachingRepository,
) CreateProjectCategoryHandler {
	return CreateProjectCategoryHandler{
		projectRepository:         projectRepository,
		projectCategoryRepository: projectCategoryRepository,
		cachingRepository:         cachingRepository,
	}
}

// CreateProjectCategory godoc
// @tags     Project
// @summary  Create project category
// @id       project-create-category
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload body    dto.CreateProjectCategoryRequest true "Body"
// @success  200     {object} dto.CreateProjectCategoryResponse
// @router   /api/project/{id}/category [post]
func (h CreateProjectCategoryHandler) CreateProjectCategory(ctx *appcontext.AppContext, performerID, projectID string, req dto.CreateProjectCategoryRequest) (*dto.CreateProjectCategoryResponse, error) {
	ctx.Logger().Info("new create project category request", appcontext.Fields{
		"performerID": performerID, "projectID": projectID, "name": req.Name,
	})

	ctx.Logger().Text("find project in db")
	project, err := h.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}
	if !project.IsOwner(performerID) {
		ctx.Logger().ErrorText("user is not project owner")
		return nil, apperrors.Common.NotFound
	}

	ctx.Logger().Text("count total created categories of this project")
	total, err := h.projectCategoryRepository.CountTotalByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to count total created categories of this project", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Info("total created categories of this project", appcontext.Fields{"total": total})

	if domain.IsReachedMaxCategoryPerProject(total) {
		ctx.Logger().ErrorText("project has too many categories")
		return nil, apperrors.Project.CategoryLimitExceeded
	}

	ctx.Logger().Text("create project category model")
	category, err := domain.NewProjectCategory(projectID, req.Name)
	if err != nil {
		ctx.Logger().Error("failed to create project category model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist project category in db")
	if err = h.projectCategoryRepository.Create(ctx, *category); err != nil {
		ctx.Logger().Error("failed to persist project in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectCategoriesByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done create project category request")
	return &dto.CreateProjectCategoryResponse{ID: category.ID}, nil
}
