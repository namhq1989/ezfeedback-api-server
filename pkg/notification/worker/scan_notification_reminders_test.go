package worker_test

import (
	"context"
	"testing"

	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type scanNotificationRemindersTestSuite struct {
	suite.Suite
	handler                            worker.ScanNotificationRemindersHandler
	mockCtrl                           *gomock.Controller
	mockNotificationReminderRepository *mocknotification.MockNotificationReminderRepository
	mockQueueRepository                *mocknotification.MockQueueRepository
}

func (s *scanNotificationRemindersTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *scanNotificationRemindersTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationReminderRepository = mocknotification.NewMockNotificationReminderRepository(s.mockCtrl)
	s.mockQueueRepository = mocknotification.NewMockQueueRepository(s.mockCtrl)

	s.handler = worker.NewScanNotificationRemindersHandler(s.mockNotificationReminderRepository, s.mockQueueRepository)
}

func (s *scanNotificationRemindersTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *scanNotificationRemindersTestSuite) Test_1_Success_HasReminders() {
	// mock
	s.mockNotificationReminderRepository.EXPECT().
		FindAllExisting(gomock.Any()).
		Return([]domain.NotificationReminder{
			{ID: uuid.New()},
		}, nil)

	s.mockQueueRepository.EXPECT().
		ProcessNotificationReminder(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.ScanNotificationReminders(ctx, domain.QueueScanNotificationRemindersPayload{})

	assert.Nil(s.T(), err)
}

func (s *scanNotificationRemindersTestSuite) Test_1_Success_NoReminders() {
	// mock
	s.mockNotificationReminderRepository.EXPECT().
		FindAllExisting(gomock.Any()).
		Return([]domain.NotificationReminder{}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.ScanNotificationReminders(ctx, domain.QueueScanNotificationRemindersPayload{})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestScanNotificationRemindersTestSuite(t *testing.T) {
	suite.Run(t, new(scanNotificationRemindersTestSuite))
}
