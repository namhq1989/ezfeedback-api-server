package command_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type createProjectCampaignTestSuite struct {
	suite.Suite
	handler                       command.CreateProjectCampaignHandler
	mockCtrl                      *gomock.Controller
	mockProjectRepository         *mockproject.MockProjectRepository
	mockProjectCampaignRepository *mockproject.MockProjectCampaignRepository
	mockCachingRepository         *mockproject.MockCachingRepository
}

func (s *createProjectCampaignTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectCampaignRepository = mockproject.NewMockProjectCampaignRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.handler = command.NewCreateProjectCampaignHandler(s.mockProjectRepository, s.mockProjectCampaignRepository, s.mockCachingRepository)
}

func (s *createProjectCampaignTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createProjectCampaignTestSuite) Test_1_Success() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)
	s.mockProjectCampaignRepository.EXPECT().
		CountTotalByProjectIDAndCampaignType(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(0), nil)
	s.mockProjectCampaignRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)
	s.mockCachingRepository.EXPECT().
		DeleteProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)
	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCampaign(ctx, performerID, projectID, dto.CreateProjectCampaignRequest{
		Name:           "Campaign 1",
		Description:    "Campaign description",
		CampaignType:   domain.ProjectCampaignTypeFeedback.String(),
		WidgetPosition: "bottom-right",
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.ID)
}

func (s *createProjectCampaignTestSuite) Test_2_Fail_ProjectNotFound() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCampaign(ctx, uuid.New(), uuid.New(), dto.CreateProjectCampaignRequest{
		Name:         "Campaign 1",
		Description:  "Campaign description",
		CampaignType: domain.ProjectCampaignTypeFeedback.String(),
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *createProjectCampaignTestSuite) Test_2_Fail_NotOwner() {
	projectID := uuid.New()
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: uuid.New(), Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCampaign(ctx, uuid.New(), projectID, dto.CreateProjectCampaignRequest{
		Name:         "Campaign 1",
		Description:  "Campaign description",
		CampaignType: domain.ProjectCampaignTypeFeedback.String(),
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.NotFound, err)
}

func (s *createProjectCampaignTestSuite) Test_3_Fail_InvalidName() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCampaign(ctx, performerID, projectID, dto.CreateProjectCampaignRequest{
		Name:         "a",
		Description:  "Campaign description",
		CampaignType: domain.ProjectCampaignTypeFeedback.String(),
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidName, err)
}

func (s *createProjectCampaignTestSuite) Test_4_Fail_InvalidCampaignType() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCampaign(ctx, performerID, projectID, dto.CreateProjectCampaignRequest{
		Name:         "Campaign 1",
		Description:  "Campaign description",
		CampaignType: "invalid_type",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaignType, err)
}

//
// END OF CASES
//

func TestCreateProjectCampaignTestSuite(t *testing.T) {
	suite.Run(t, new(createProjectCampaignTestSuite))
}
