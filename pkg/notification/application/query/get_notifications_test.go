package query_test

import (
	"context"
	"testing"

	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getNotificationsTestSuite struct {
	suite.Suite
	handler                    query.GetNotificationsHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *getNotificationsTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)
	s.handler = query.NewGetNotificationsHandler(
		s.mockNotificationRepository,
	)
}

func (s *getNotificationsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getNotificationsTestSuite) Test_1_Success() {
	s.mockNotificationRepository.EXPECT().
		FindWithFilter(gomock.Any(), gomock.Any()).
		Return([]domain.Notification{
			{ID: uuid.New()},
		}, nil)

	s.mockNotificationRepository.EXPECT().
		GenerateNewFeedbackContent(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("content").
		AnyTimes()

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetNotifications(ctx, uuid.New(), dto.GetNotificationsRequest{
		Page: 0,
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetNotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(getNotificationsTestSuite))
}
