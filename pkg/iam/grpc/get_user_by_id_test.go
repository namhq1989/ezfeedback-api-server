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

type getUserByIDTestSuite struct {
	suite.Suite
	handler     grpc.GetUserByIDHandler
	mockCtrl    *gomock.Controller
	mockService *mockiam.MockService
}

func (s *getUserByIDTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getUserByIDTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockiam.NewMockService(s.mockCtrl)

	s.handler = grpc.NewGetUserByIDHandler(s.mockService)
}

func (s *getUserByIDTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getUserByIDTestSuite) Test_1_Success() {
	// mock
	s.mockService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil)
	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetUserById(ctx, &iampb.GetUserByIdRequest{
		TraceId: "trace-id",
		UserId:  uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetUserByIDTestSuite(t *testing.T) {
	suite.Run(t, new(getUserByIDTestSuite))
}
