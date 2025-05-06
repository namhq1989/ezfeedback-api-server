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

type getProjectCampaignsTestSuite struct {
	suite.Suite
	handler     grpc.GetProjectCampaignsHandler
	mockCtrl    *gomock.Controller
	mockService *mockproject.MockService
}

func (s *getProjectCampaignsTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getProjectCampaignsTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = grpc.NewGetProjectCampaignsHandler(s.mockService)
}

func (s *getProjectCampaignsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectCampaignsTestSuite) Test_1_Success() {
	// mock
	s.mockService.EXPECT().
		GetProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCampaign{{ID: uuid.New()}}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetProjectCampaigns(ctx, &projectpb.GetProjectCampaignsRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.GetCampaigns()))
}

//
// END OF CASES
//

func TestGetProjectCampaignsTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectCampaignsTestSuite))
}
