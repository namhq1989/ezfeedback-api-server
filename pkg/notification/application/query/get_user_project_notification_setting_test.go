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

type getUserProjectNotificationSettingTestSuite struct {
	suite.Suite
	handler     query.GetUserProjectNotificationSettingHandler
	mockCtrl    *gomock.Controller
	mockService *mocknotification.MockService
}

func (s *getUserProjectNotificationSettingTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mocknotification.NewMockService(s.mockCtrl)
	s.handler = query.NewGetUserProjectNotificationSettingHandler(
		s.mockService,
	)
}

func (s *getUserProjectNotificationSettingTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getUserProjectNotificationSettingTestSuite) Test_1_Success() {
	s.mockService.EXPECT().
		GetUserProjectNotificationSetting(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.UserProjectNotificationSetting{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetUserProjectNotificationSetting(ctx, uuid.New(), dto.GetUserProjectNotificationSettingRequest{
		ProjectID: uuid.New(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetUserProjectNotificationSettingTestSuite(t *testing.T) {
	suite.Run(t, new(getUserProjectNotificationSettingTestSuite))
}
