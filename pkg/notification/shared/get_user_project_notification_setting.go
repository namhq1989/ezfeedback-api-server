package shared

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string) (*domain.UserProjectNotificationSetting, error) {
	ctx.Logger().Info("[service] get user project notification setting", appcontext.Fields{"userID": userID, "projectID": projectID})

	ctx.Logger().Text("find user project notification setting in caching")
	setting, err := s.cachingRepository.GetUserProjectNotificationSetting(ctx, userID, projectID)
	if setting != nil {
		ctx.Logger().Text("user project notification setting found in caching, return")
		return setting, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find user project notification setting in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("user project notification setting not found in caching, find in db")
	setting, err = s.userProjectNotificationSettingRepository.FindByUserIDAndProjectID(ctx, userID, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find user project notification setting in db", err, appcontext.Fields{})
		return nil, err
	}
	if setting == nil {
		ctx.Logger().Text("user project notification setting not found in db, create a new one")
		setting, err = domain.NewUserProjectNotificationSetting(userID, projectID)
		if err != nil {
			ctx.Logger().Error("failed to create user project notification setting", err, appcontext.Fields{})
			return nil, err
		}

		ctx.Logger().Text("persist user project notification setting in db")
		if err = s.userProjectNotificationSettingRepository.Create(ctx, *setting); err != nil {
			ctx.Logger().Error("failed to persist user project notification setting in db", err, appcontext.Fields{})
			return nil, err
		}
	}

	ctx.Logger().Text("set user project notification setting in caching")
	if err = s.cachingRepository.SetUserProjectNotificationSetting(ctx, userID, projectID, *setting); err != nil {
		ctx.Logger().Error("failed to set user project notification setting in caching", err, appcontext.Fields{})
	}
	return setting, nil
}
