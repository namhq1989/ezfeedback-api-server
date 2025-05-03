package infrastructure

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CachingRepository struct {
	caching caching.Operations

	domain                                    string
	userProjectNotificationSettingCachingTime time.Duration
}

func NewCachingRepository(caching *caching.Caching, isEnvRelease bool) CachingRepository {
	if isEnvRelease {
		return CachingRepository{
			caching: caching,
			domain:  "notification",
			userProjectNotificationSettingCachingTime: 12 * time.Hour,
		}
	} else {
		cachingTime := 1 * time.Minute

		return CachingRepository{
			caching: caching,
			domain:  "notification",
			userProjectNotificationSettingCachingTime: cachingTime,
		}
	}
}

//
// GET USER PROJECT NOTIFICATION SETTING
//

func (r CachingRepository) GetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string) (*domain.UserProjectNotificationSetting, error) {
	key := r.generateUserProjectNotificationSettingKey(userID, projectID)

	dataStr, err := r.caching.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result *domain.UserProjectNotificationSetting
	if err = json.Unmarshal([]byte(dataStr), &result); err != nil {
		return nil, nil
	}

	return result, nil
}

func (r CachingRepository) SetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string, setting domain.UserProjectNotificationSetting) error {
	key := r.generateUserProjectNotificationSettingKey(userID, projectID)
	r.caching.SetTTL(ctx, key, setting, r.userProjectNotificationSettingCachingTime)
	return nil
}

func (r CachingRepository) DeleteProjectByID(ctx *appcontext.AppContext, userID, projectID string) error {
	key := r.generateUserProjectNotificationSettingKey(userID, projectID)
	_, err := r.caching.Del(ctx, key)
	return err
}

func (r CachingRepository) generateUserProjectNotificationSettingKey(userID, projectID string) string {
	return r.caching.GenerateKey(r.domain, fmt.Sprintf("user:%s:project:%s:setting", userID, projectID))
}
