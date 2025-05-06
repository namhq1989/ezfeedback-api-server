package worker_test

import (
	"context"
	"testing"

	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type cleanupStaleTestSuite struct {
	suite.Suite
	handler                    worker.CleanupStaleNotificationsHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *cleanupStaleTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *cleanupStaleTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)

	s.handler = worker.NewCleanupStaleNotificationsHandler(s.mockNotificationRepository)
}

func (s *cleanupStaleTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *cleanupStaleTestSuite) Test_1_Success() {
	// mock
	s.mockNotificationRepository.EXPECT().
		CleanupStale(gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.CleanupStaleNotifications(ctx, domain.QueueCleanupStaleNotificationsPayload{})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestCleanupStaleTestSuite(t *testing.T) {
	suite.Run(t, new(cleanupStaleTestSuite))
}
