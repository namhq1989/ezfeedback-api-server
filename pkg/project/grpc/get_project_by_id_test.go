package grpc_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getProjectByIDTestSuite struct {
	suite.Suite
	handler     grpc.GetProjectByIDHandler
	mockCtrl    *gomock.Controller
	mockService *mockproject.MockService
}

func (s *getProjectByIDTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getProjectByIDTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = grpc.NewGetProjectByIDHandler(s.mockService)
}

func (s *getProjectByIDTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectByIDTestSuite) Test_1_Success() {
	// mock
	s.mockService.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.Project{
			ID:     uuid.New(),
			Status: domain.StatusActive,
		}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetProjectById(ctx, &projectpb.GetProjectByIdRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetProjectByIDTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectByIDTestSuite))
}
