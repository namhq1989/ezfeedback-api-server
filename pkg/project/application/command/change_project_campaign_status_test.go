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

type changeProjectCampaignStatusTestSuite struct {
	suite.Suite
	handler                       command.ChangeProjectCampaignStatusHandler
	mockCtrl                      *gomock.Controller
	mockProjectCampaignRepository *mockproject.MockProjectCampaignRepository
	mockCachingRepository         *mockproject.MockCachingRepository
	mockService                   *mockproject.MockService
}

func (s *changeProjectCampaignStatusTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCampaignRepository = mockproject.NewMockProjectCampaignRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)
	s.handler = command.NewChangeProjectCampaignStatusHandler(s.mockProjectCampaignRepository, s.mockCachingRepository, s.mockService)
}

func (s *changeProjectCampaignStatusTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *changeProjectCampaignStatusTestSuite) Test_1_Success() {
	projectID := uuid.New()
	campaignID := uuid.New()
	performerID := uuid.New()
	status := domain.StatusInactive.String()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, dto.ChangeProjectCampaignStatusRequest{
		Status: status,
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectCampaignStatusTestSuite) Test_1_Success_StatusNotChanged() {
	projectID := uuid.New()
	campaignID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusInactive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.ProjectNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, uuid.New(), "invalid-id", uuid.New(), dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_InvalidCampaignID() {
	projectID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCampaign)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, "invalid-campaign-id", dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaign, err)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_NotOwner() {
	projectID := uuid.New()
	campaignID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Common.NotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.NotFound, err)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_InvalidStatus() {
	projectID := uuid.New()
	campaignID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, dto.ChangeProjectCampaignStatusRequest{
		Status: "invalid-status",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.InvalidStatus, err)
}

func TestChangeProjectCampaignStatusTestSuite(t *testing.T) {
	suite.Run(t, new(changeProjectCampaignStatusTestSuite))
}
