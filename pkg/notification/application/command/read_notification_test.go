package command_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type readNotificationTestSuite struct {
	suite.Suite
	handler                    command.ReadNotificationHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *readNotificationTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)
	s.handler = command.NewReadNotificationHandler(s.mockNotificationRepository)
}

func (s *readNotificationTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *readNotificationTestSuite) Test_1_Success() {
	var (
		notificationID = uuid.New()
		performerID    = uuid.New()
	)

	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Notification{ID: notificationID, UserID: performerID}, nil)

	s.mockNotificationRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ReadNotification(ctx, performerID, notificationID, dto.ReadNotificationRequest{})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *readNotificationTestSuite) Test_2_Fail_NotificationNotFound() {
	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ReadNotification(ctx, uuid.New(), uuid.New(), dto.ReadNotificationRequest{})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Notification.NotificationNotFound, err)
}

func (s *readNotificationTestSuite) Test_2_Fail_NotOwner() {
	notificationID := uuid.New()

	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Notification{ID: notificationID, UserID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ReadNotification(ctx, uuid.New(), notificationID, dto.ReadNotificationRequest{})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Notification.NotificationNotFound, err)
}

//
// END OF CASES
//

func TestReadNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(readNotificationTestSuite))
}
