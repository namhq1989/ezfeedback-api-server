package shared_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/shared"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getUserByIDTestSuite struct {
	suite.Suite
	service               shared.Service
	mockCtrl              *gomock.Controller
	mockUserRepository    *mockiam.MockUserRepository
	mockCachingRepository *mockiam.MockCachingRepository
}

func (s *getUserByIDTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getUserByIDTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockUserRepository = mockiam.NewMockUserRepository(s.mockCtrl)
	s.mockCachingRepository = mockiam.NewMockCachingRepository(s.mockCtrl)

	s.service = shared.NewService(s.mockUserRepository, s.mockCachingRepository)
}

func (s *getUserByIDTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getUserByIDTestSuite) Test_1_Success_UserFoundInCaching() {
	user := &domain.User{ID: uuid.New()}

	s.mockCachingRepository.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(user, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.service.GetUserByID(ctx, uuid.New())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), user, resp)
}

func (s *getUserByIDTestSuite) Test_1_Success_UserNotFoundInCaching() {
	user := &domain.User{ID: uuid.New()}

	s.mockCachingRepository.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockUserRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(user, nil)

	s.mockCachingRepository.EXPECT().
		SetUserByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.service.GetUserByID(ctx, uuid.New())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), user, resp)
}

func (s *getUserByIDTestSuite) Test_2_Fail_UserNotFound() {
	s.mockCachingRepository.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockUserRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.User.UserNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.service.GetUserByID(ctx, uuid.New())
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.User.UserNotFound, err)
}

//
// END OF CASES
//

func TestGetUserByIDTestSuite(t *testing.T) {
	suite.Run(t, new(getUserByIDTestSuite))
}
