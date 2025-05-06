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

type deleteNotificationTestSuite struct {
	suite.Suite
	handler                    command.DeleteNotificationHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *deleteNotificationTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)
	s.handler = command.NewDeleteNotificationHandler(s.mockNotificationRepository)
}

func (s *deleteNotificationTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *deleteNotificationTestSuite) Test_1_Success() {
	var (
		notificationID = uuid.New()
		performerID    = uuid.New()
	)

	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Notification{ID: notificationID, UserID: performerID}, nil)

	s.mockNotificationRepository.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.DeleteNotification(ctx, performerID, notificationID, dto.DeleteNotificationRequest{})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *deleteNotificationTestSuite) Test_2_Fail_NotificationNotFound() {
	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.DeleteNotification(ctx, uuid.New(), uuid.New(), dto.DeleteNotificationRequest{})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Notification.NotificationNotFound, err)
}

func (s *deleteNotificationTestSuite) Test_2_Fail_NotOwner() {
	notificationID := uuid.New()

	s.mockNotificationRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Notification{ID: notificationID, UserID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.DeleteNotification(ctx, uuid.New(), notificationID, dto.DeleteNotificationRequest{})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Notification.NotificationNotFound, err)
}

//
// END OF CASES
//

func TestDeleteNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(deleteNotificationTestSuite))
}
