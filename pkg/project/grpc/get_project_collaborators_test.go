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

type getProjectCollaboratorsTestSuite struct {
	suite.Suite
	handler     grpc.GetProjectCollaboratorsHandler
	mockCtrl    *gomock.Controller
	mockService *mockproject.MockService
}

func (s *getProjectCollaboratorsTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getProjectCollaboratorsTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = grpc.NewGetProjectCollaboratorsHandler(s.mockService)
}

func (s *getProjectCollaboratorsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectCollaboratorsTestSuite) Test_1_Success() {
	// mock
	s.mockService.EXPECT().
		GetProjectCollaboratorsByProjectID(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCollaborator{{ID: uuid.New()}}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetProjectCollaborators(ctx, &projectpb.GetProjectCollaboratorsRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.GetCollaborators()))
}

//
// END OF CASES
//

func TestGetProjectCollaboratorsTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectCollaboratorsTestSuite))
}
