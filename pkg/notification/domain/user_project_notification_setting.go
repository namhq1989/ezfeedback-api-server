package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type UserProjectNotificationSettingRepository interface {
	Create(ctx *appcontext.AppContext, setting UserProjectNotificationSetting) error
	Update(ctx *appcontext.AppContext, setting UserProjectNotificationSetting) error
	FindByUserIDAndProjectID(ctx *appcontext.AppContext, userID, projectID string) (*UserProjectNotificationSetting, error)
}

type UserProjectNotificationSetting struct {
	ID                   string
	UserID               string
	ProjectID            string
	ReceiveNewFeedback   bool
	ReceiveDailySummary  bool
	ReceiveWeeklySummary bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewUserProjectNotificationSetting(userID, projectID string) (*UserProjectNotificationSetting, error) {
	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var (
		now = manipulation.NowUTC()
	)

	return &UserProjectNotificationSetting{
		ID:                   uuid.New(),
		UserID:               userID,
		ProjectID:            projectID,
		ReceiveNewFeedback:   true,
		ReceiveDailySummary:  true,
		ReceiveWeeklySummary: true,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func (s *UserProjectNotificationSetting) SetReceiveNewFeedback(value bool) {
	s.ReceiveNewFeedback = value
	s.SetUpdatedAt()
}

func (s *UserProjectNotificationSetting) SetReceiveDailySummary(value bool) {
	s.ReceiveDailySummary = value
	s.SetUpdatedAt()
}

func (s *UserProjectNotificationSetting) SetReceiveWeeklySummary(value bool) {
	s.ReceiveWeeklySummary = value
	s.SetUpdatedAt()
}

func (s *UserProjectNotificationSetting) SetUpdatedAt() {
	s.UpdatedAt = manipulation.NowUTC()
}
