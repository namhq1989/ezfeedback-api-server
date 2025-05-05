package worker_test

import (
	"context"
	"testing"
	"time"

	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type processNotificationReminderTestSuite struct {
	suite.Suite
	handler                            worker.ProcessNotificationReminderHandler
	mockCtrl                           *gomock.Controller
	mockNotificationReminderRepository *mocknotification.MockNotificationReminderRepository
	mockMailerRepository               *mocknotification.MockMailerRepository
	mockFeedbackHub                    *mocknotification.MockFeedbackHub
	mockProjectHub                     *mocknotification.MockProjectHub
}

func (s *processNotificationReminderTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *processNotificationReminderTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationReminderRepository = mocknotification.NewMockNotificationReminderRepository(s.mockCtrl)
	s.mockMailerRepository = mocknotification.NewMockMailerRepository(s.mockCtrl)
	s.mockFeedbackHub = mocknotification.NewMockFeedbackHub(s.mockCtrl)
	s.mockProjectHub = mocknotification.NewMockProjectHub(s.mockCtrl)

	s.handler = worker.NewProcessNotificationReminderHandler(
		s.mockNotificationReminderRepository,
		s.mockMailerRepository,
		s.mockFeedbackHub,
		s.mockProjectHub,
	)
}

func (s *processNotificationReminderTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *processNotificationReminderTestSuite) Test_1_Success_HasReminders() {
	// mock
	s.mockFeedbackHub.EXPECT().
		GetProjectStatsForNotificationReminder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(10), []domain.Feedback{
			{ID: uuid.New()},
		}, nil)

	s.mockProjectHub.EXPECT().
		GetProjectCollaborators(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCollaborator{
			{
				ID: uuid.New(), User: domain.User{
					ID:    uuid.New(),
					Name:  "test",
					Email: "test",
				},
			},
		}, nil)

	s.mockMailerRepository.EXPECT().
		SendNewFeedbackEmail(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.ProcessNotificationReminder(ctx, domain.QueueProcessNotificationReminderPayload{
		Reminder: domain.NotificationReminder{
			ID:        uuid.New(),
			ProjectID: uuid.New(),
			CreatedAt: time.Now(),
		},
	})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestProcessNotificationReminderTestSuite(t *testing.T) {
	suite.Run(t, new(processNotificationReminderTestSuite))
}
