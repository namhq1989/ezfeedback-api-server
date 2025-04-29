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

type updateProjectCampaignTestSuite struct {
	suite.Suite
	handler                       command.UpdateProjectCampaignHandler
	mockCtrl                      *gomock.Controller
	mockProjectCampaignRepository *mockproject.MockProjectCampaignRepository
	mockCachingRepository         *mockproject.MockCachingRepository
	mockService                   *mockproject.MockService
}

func (s *updateProjectCampaignTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCampaignRepository = mockproject.NewMockProjectCampaignRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)
	s.handler = command.NewUpdateProjectCampaignHandler(s.mockProjectCampaignRepository, s.mockCachingRepository, s.mockService)
}

func (s *updateProjectCampaignTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *updateProjectCampaignTestSuite) Test_1_Success() {
	var (
		projectID   = uuid.New()
		campaignID  = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Name: "OldName"}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCampaign(ctx, performerID, projectID, campaignID, dto.UpdateProjectCampaignRequest{
		Name: "UpdatedName",
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *updateProjectCampaignTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.ProjectNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCampaign(ctx, uuid.New(), "invalid-id", uuid.New(), dto.UpdateProjectCampaignRequest{
		Name: "Test",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *updateProjectCampaignTestSuite) Test_2_Fail_InvalidCampaignD() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCampaign)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCampaign(ctx, performerID, projectID, "invalid-campaign-id", dto.UpdateProjectCampaignRequest{
		Name: "Test",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaign, err)
}

func (s *updateProjectCampaignTestSuite) Test_2_Fail_InvalidName() {
	var (
		projectID   = uuid.New()
		campaignID  = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Name: "OldName"}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCampaign(ctx, performerID, projectID, campaignID, dto.UpdateProjectCampaignRequest{
		Name: "",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.InvalidName, err)
}

func TestUpdateProjectCampaignTestSuite(t *testing.T) {
	suite.Run(t, new(updateProjectCampaignTestSuite))
}
