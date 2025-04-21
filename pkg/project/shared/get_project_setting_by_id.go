package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectSettingByID(ctx *appcontext.AppContext, projectID string) (*domain.ProjectSetting, error) {
	ctx.Logger().Info("[service] get project setting by id", appcontext.Fields{"projectID": projectID})

	ctx.Logger().Text("find project setting in caching")
	setting, err := s.cachingRepository.GetProjectSettingByProjectID(ctx, projectID)
	if setting != nil {
		ctx.Logger().Text("project setting found in caching, return")
		return setting, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project setting in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project setting not found in caching, find in db")
	setting, err = s.projectSettingRepository.FindByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project setting in db", err, appcontext.Fields{})
		return nil, err
	}
	if setting == nil {
		ctx.Logger().ErrorText("project setting not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("set project setting in caching")
	if err = s.cachingRepository.SetProjectSettingByProjectID(ctx, projectID, *setting); err != nil {
		ctx.Logger().Error("failed to set project setting in caching", err, appcontext.Fields{})
	}
	return setting, nil
}
