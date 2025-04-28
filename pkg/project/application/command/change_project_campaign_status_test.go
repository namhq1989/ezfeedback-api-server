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
	mockProjectRepository         *mockproject.MockProjectRepository
	mockProjectCampaignRepository *mockproject.MockProjectCampaignRepository
	mockCachingRepository         *mockproject.MockCachingRepository
}

func (s *changeProjectCampaignStatusTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectCampaignRepository = mockproject.NewMockProjectCampaignRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.handler = command.NewChangeProjectCampaignStatusHandler(s.mockProjectRepository, s.mockProjectCampaignRepository, s.mockCachingRepository)
}

func (s *changeProjectCampaignStatusTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *changeProjectCampaignStatusTestSuite) Test_1_Success() {
	projectID := uuid.New()
	campaignID := uuid.New()
	performerID := uuid.New()
	status := domain.StatusInactive.String()

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
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

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusInactive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, performerID, projectID, campaignID, dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidProjectID)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCampaignStatus(ctx, uuid.New(), "invalid-id", uuid.New(), dto.ChangeProjectCampaignStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.InvalidProjectID, err)
}

func (s *changeProjectCampaignStatusTestSuite) Test_2_Fail_InvalidCampaignID() {
	projectID := uuid.New()
	performerID := uuid.New()

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
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
	otherUserID := uuid.New()

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: otherUserID}, nil)

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

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID}, nil)

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
