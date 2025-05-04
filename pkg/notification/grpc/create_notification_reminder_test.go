package grpc_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type createNotificationReminderTestSuite struct {
	suite.Suite
	handler                            grpc.CreateNotificationReminderHandler
	mockCtrl                           *gomock.Controller
	mockNotificationReminderRepository *mocknotification.MockNotificationReminderRepository
}

func (s *createNotificationReminderTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *createNotificationReminderTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationReminderRepository = mocknotification.NewMockNotificationReminderRepository(s.mockCtrl)

	s.handler = grpc.NewCreateNotificationReminderHandler(s.mockNotificationReminderRepository)
}

func (s *createNotificationReminderTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createNotificationReminderTestSuite) Test_1_Success_ReminderExisted() {
	// mock
	s.mockNotificationReminderRepository.EXPECT().
		FindByProjectID(gomock.Any(), gomock.Any()).
		Return(&domain.NotificationReminder{
			ID: uuid.New(),
		}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.CreateNotificationReminder(ctx, &notificationpb.CreateNotificationReminderRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *createNotificationReminderTestSuite) Test_1_Success_NewReminder() {
	// mock
	s.mockNotificationReminderRepository.EXPECT().
		FindByProjectID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockNotificationReminderRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.CreateNotificationReminder(ctx, &notificationpb.CreateNotificationReminderRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestCreateNotificationReminderTestSuite(t *testing.T) {
	suite.Run(t, new(createNotificationReminderTestSuite))
}
