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

type getProjectCampaignByIDTestSuite struct {
	suite.Suite
	handler                grpc.GetProjectCampaignByIDHandler
	mockCtrl               *gomock.Controller
	mockProjectCampaignHub *mockproject.MockProjectCampaignHub
}

func (s *getProjectCampaignByIDTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getProjectCampaignByIDTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCampaignHub = mockproject.NewMockProjectCampaignHub(s.mockCtrl)

	s.handler = grpc.NewGetProjectCampaignByIDHandler(s.mockProjectCampaignHub)
}

func (s *getProjectCampaignByIDTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectCampaignByIDTestSuite) Test_1_Success() {
	// mock
	s.mockProjectCampaignHub.EXPECT().
		FindProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaignHubData{
			Project: domain.Project{
				ID:     uuid.New(),
				Status: domain.StatusActive,
			},
			ProjectCampaign: domain.ProjectCampaign{
				ID:           uuid.New(),
				CampaignType: domain.ProjectCampaignTypeFeedback,
				Status:       domain.StatusActive,
			},
			ProjectSetting: domain.ProjectSetting{
				Domain: "domain.com",
			},
		}, nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetProjectCampaignByID(ctx, &projectpb.GetProjectCampaignByIdRequest{
		TraceId:    "trace-id",
		CampaignId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetProjectCampaignByIDTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectCampaignByIDTestSuite))
}
