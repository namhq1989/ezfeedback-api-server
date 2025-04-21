package command_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type updateMeTestSuite struct {
	suite.Suite
	handler               command.UpdateMeHandler
	mockCtrl              *gomock.Controller
	mockUserRepository    *mockiam.MockUserRepository
	mockCachingRepository *mockiam.MockCachingRepository
	mockService           *mockiam.MockService
}

func (s *updateMeTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockUserRepository = mockiam.NewMockUserRepository(s.mockCtrl)
	s.mockCachingRepository = mockiam.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockiam.NewMockService(s.mockCtrl)
	s.handler = command.NewUpdateMeHandler(s.mockUserRepository, s.mockCachingRepository, s.mockService)
}

func (s *updateMeTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *updateMeTestSuite) Test_1_Success() {
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil)

	s.mockUserRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteUserByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateMe(ctx, uuid.New(), dto.UpdateMeRequest{Name: "New Name"})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *updateMeTestSuite) Test_2_Fail_InvalidName() {
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateMe(ctx, uuid.New(), dto.UpdateMeRequest{Name: ""})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidName, err)
}

func (s *updateMeTestSuite) Test_3_Fail_UserNotFound() {
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.User.UserNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateMe(ctx, uuid.New(), dto.UpdateMeRequest{Name: "New Name"})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.User.UserNotFound, err)
}

//
// END OF CASES
//

func TestUpdateMeTestSuite(t *testing.T) {
	suite.Run(t, new(updateMeTestSuite))
}
