package grpc_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getUsersByIDsTestSuite struct {
	suite.Suite
	handler     grpc.GetUsersByIDsHandler
	mockCtrl    *gomock.Controller
	mockService *mockiam.MockService
}

func (s *getUsersByIDsTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getUsersByIDsTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockiam.NewMockService(s.mockCtrl)

	s.handler = grpc.NewGetUsersByIDsHandler(s.mockService)
}

func (s *getUsersByIDsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getUsersByIDsTestSuite) Test_1_Success() {
	// mock
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil).
		AnyTimes()

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetUsersByIds(ctx, &iampb.GetUsersByIdsRequest{
		TraceId: "trace-id",
		UserIds: []string{uuid.New()},
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.GetUsers()))
}

//
// END OF CASES
//

func TestGetUsersByIDsTestSuite(t *testing.T) {
	suite.Run(t, new(getUsersByIDsTestSuite))
}
