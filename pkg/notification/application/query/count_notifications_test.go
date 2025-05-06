package query_test

import (
	"context"
	"testing"

	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type countNotificationsTestSuite struct {
	suite.Suite
	handler                    query.CountNotificationsHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *countNotificationsTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)
	s.handler = query.NewCountNotificationsHandler(s.mockNotificationRepository)
}

func (s *countNotificationsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *countNotificationsTestSuite) Test_1_Success() {
	s.mockNotificationRepository.EXPECT().
		CountWithFilter(gomock.Any(), gomock.Any()).
		Return(int64(10), nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CountNotifications(ctx, uuid.New(), dto.CountNotificationsRequest{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), int64(10), resp.Total)
}

//
// END OF CASES
//

func TestCountNotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(countNotificationsTestSuite))
}
