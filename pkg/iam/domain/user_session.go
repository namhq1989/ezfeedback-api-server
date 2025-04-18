package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

var (
	sessionTTL = 30 * 24 * time.Hour
)

type UserSession struct {
	ID           string
	UserID       string
	DeviceID     string
	RefreshToken string
	DeviceInfo   UserDeviceInfo
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUserSession(userID, deviceID, refreshToken, ip, userAgent string) (*UserSession, error) {
	var (
		now = manipulation.NowUTC()
	)

	s := &UserSession{
		ID:        uuid.New(),
		ExpiresAt: now.Add(sessionTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := s.SetDeviceInfo(ip, userAgent); err != nil {
		return nil, err
	}
	if err := s.SetRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *UserSession) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}
	s.UserID = userID
	return nil
}

func (s *UserSession) SetDeviceID(deviceID string) error {
	if deviceID == "" {
		return apperrors.User.InvalidDeviceID
	}
	s.DeviceID = deviceID
	return nil
}

func (s *UserSession) SetRefreshToken(token string) error {
	if token == "" {
		return apperrors.Auth.InvalidRefreshToken
	}
	s.RefreshToken = token
	return nil
}

func (s *UserSession) SetDeviceInfo(ip, userAgent string) error {
	os := manipulation.ExtractOSFromUserAgent(userAgent)
	s.DeviceInfo = UserDeviceInfo{
		OS:        os,
		UserAgent: userAgent,
		Ip:        ip,
	}

	return nil
}

func (s *UserSession) IsExpired() bool {
	return s.ExpiresAt.Before(manipulation.NowUTC())
}
