package shared

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectCategoriesByProjectID(ctx *appcontext.AppContext, projectID string, status domain.Status) ([]domain.ProjectCategory, error) {
	ctx.Logger().Info("[service] get project categories by project id", appcontext.Fields{"projectID": projectID, "status": status.String()})

	ctx.Logger().Text("find project categories in caching")
	categories, err := s.cachingRepository.GetProjectCategoriesByProjectID(ctx, projectID)
	if categories != nil {
		ctx.Logger().Text("project categories found in caching, return")
		return s.filterProjectCategories(ctx, categories, status), nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project categories in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project categories not found in caching, find in db")
	categories, err = s.projectCategoryRepository.FindByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project categories in db", err, appcontext.Fields{})
		return nil, err
	}
	if categories == nil || len(categories) == 0 {
		ctx.Logger().Text("project categories not found")
		return make([]domain.ProjectCategory, 0), nil
	}

	ctx.Logger().Text("set project categories in caching")
	if err = s.cachingRepository.SetProjectCategoriesByProjectID(ctx, projectID, categories); err != nil {
		ctx.Logger().Error("failed to set project categories in caching", err, appcontext.Fields{})
	}
	return s.filterProjectCategories(ctx, categories, status), nil
}

func (s Service) filterProjectCategories(_ *appcontext.AppContext, categories []domain.ProjectCategory, status domain.Status) []domain.ProjectCategory {
	if !status.IsValid() {
		return categories
	}

	var result []domain.ProjectCategory
	for _, category := range categories {
		if category.Status == status {
			result = append(result, category)
		}
	}
	return result
}
