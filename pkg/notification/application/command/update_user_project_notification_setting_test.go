package command_test

import (
	"context"
	"testing"

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

type updateUserProjectNotificationSettingTestSuite struct {
	suite.Suite
	handler                                      command.UpdateUserProjectNotificationSettingHandler
	mockCtrl                                     *gomock.Controller
	mockUserProjectNotificationSettingRepository *mocknotification.MockUserProjectNotificationSettingRepository
	mockCachingRepository                        *mocknotification.MockCachingRepository
	mockService                                  *mocknotification.MockService
}

func (s *updateUserProjectNotificationSettingTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockUserProjectNotificationSettingRepository = mocknotification.NewMockUserProjectNotificationSettingRepository(s.mockCtrl)
	s.mockCachingRepository = mocknotification.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mocknotification.NewMockService(s.mockCtrl)
	s.handler = command.NewUpdateUserProjectNotificationSettingHandler(
		s.mockUserProjectNotificationSettingRepository,
		s.mockCachingRepository,
		s.mockService,
	)
}

func (s *updateUserProjectNotificationSettingTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *updateUserProjectNotificationSettingTestSuite) Test_1_Success() {
	s.mockService.EXPECT().
		GetUserProjectNotificationSetting(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.UserProjectNotificationSetting{ID: uuid.New()}, nil)

	s.mockUserProjectNotificationSettingRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		SetUserProjectNotificationSetting(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateUserProjectNotificationSetting(ctx, uuid.New(), dto.UpdateUserProjectNotificationSettingRequest{
		ProjectID: uuid.New(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestUpdateUserProjectNotificationSettingTestSuite(t *testing.T) {
	suite.Run(t, new(updateUserProjectNotificationSettingTestSuite))
}
