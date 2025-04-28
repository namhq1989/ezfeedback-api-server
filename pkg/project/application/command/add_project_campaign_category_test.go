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

type addProjectCampaignCategoryTestSuite struct {
	suite.Suite
	handler                               command.AddProjectCampaignCategoryHandler
	mockCtrl                              *gomock.Controller
	mockProjectCampaignCategoryRepository *mockproject.MockProjectCampaignCategoryRepository
	mockCachingRepository                 *mockproject.MockCachingRepository
	mockService                           *mockproject.MockService
}

func (s *addProjectCampaignCategoryTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCampaignCategoryRepository = mockproject.NewMockProjectCampaignCategoryRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)
	s.handler = command.NewAddProjectCampaignCategoryHandler(s.mockProjectCampaignCategoryRepository, s.mockCachingRepository, s.mockService)
}

func (s *addProjectCampaignCategoryTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *addProjectCampaignCategoryTestSuite) Test_1_Success() {
	projectID := uuid.New()
	campaignID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockProjectCampaignCategoryRepository.EXPECT().
		FindByCampaignIDAndCategoryID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockProjectCampaignCategoryRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.AddProjectCampaignCategory(ctx, performerID, projectID, campaignID, dto.AddProjectCampaignCategoryRequest{
		CategoryID: categoryID,
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *addProjectCampaignCategoryTestSuite) Test_1_Success_InvalidCampaignID() {
	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCampaign)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.AddProjectCampaignCategory(ctx, uuid.New(), uuid.New(), "invalid-campaign-id", dto.AddProjectCampaignCategoryRequest{
		CategoryID: uuid.New(),
	})
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaign, err)
}

func (s *addProjectCampaignCategoryTestSuite) Test_1_Success_InvalidCategoryID() {
	projectID := uuid.New()
	campaignID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCampaign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{ID: campaignID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCategory)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.AddProjectCampaignCategory(ctx, performerID, projectID, campaignID, dto.AddProjectCampaignCategoryRequest{
		CategoryID: categoryID,
	})
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCategory, err)
}

func TestAddProjectCampaignCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(addProjectCampaignCategoryTestSuite))
}
