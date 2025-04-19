package mapping

import (
	"github.com/goccy/go-json"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
)

type UserSessionMapper struct{}

func (UserSessionMapper) FromModelToDomain(session model.UserSessions) (*domain.UserSession, error) {
	var result = &domain.UserSession{
		ID:           session.ID,
		UserID:       session.UserID,
		DeviceID:     session.DeviceID,
		ExpiresAt:    session.ExpiresAt,
		RefreshToken: session.RefreshToken,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}

	if session.DeviceInfo != "" {
		if err := json.Unmarshal([]byte(session.DeviceInfo), &result.DeviceInfo); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (UserSessionMapper) FromDomainToModel(session domain.UserSession) (*model.UserSessions, error) {
	var result = &model.UserSessions{
		ID:           session.ID,
		UserID:       session.UserID,
		DeviceID:     session.DeviceID,
		ExpiresAt:    session.ExpiresAt,
		RefreshToken: session.RefreshToken,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}

	deviceInfo := UserDeviceInfo{
		OS:        session.DeviceInfo.OS,
		UserAgent: session.DeviceInfo.UserAgent,
		Ip:        session.DeviceInfo.Ip,
	}
	if data, err := json.Marshal(deviceInfo); err != nil {
		return nil, err
	} else {
		deviceInfoStr := string(data)
		result.DeviceInfo = deviceInfoStr
	}

	return result, nil
}
