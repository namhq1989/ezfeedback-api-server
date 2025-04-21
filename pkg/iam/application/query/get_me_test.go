package query_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getMeTestSuite struct {
	suite.Suite
	handler     query.GetMeHandler
	mockCtrl    *gomock.Controller
	mockService *mockiam.MockService
}

func (s *getMeTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockiam.NewMockService(s.mockCtrl)
	s.handler = query.NewGetMeHandler(s.mockService)
}

func (s *getMeTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getMeTestSuite) Test_1_Success_UserFound() {
	user := domain.User{ID: uuid.New(), Email: "john@example.com", Name: "John"}

	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&user, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetMe(ctx, user.ID, dto.GetMeRequest{})
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), user.ID, resp.Me.ID)
}

func (s *getMeTestSuite) Test_2_Fail_UserNotFound() {
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.User.UserNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetMe(ctx, uuid.New(), dto.GetMeRequest{})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.User.UserNotFound, err)
}

//
// END OF CASES
//

func TestGetMeTestSuite(t *testing.T) {
	suite.Run(t, new(getMeTestSuite))
}
