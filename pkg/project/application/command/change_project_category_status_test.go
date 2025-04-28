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

type changeProjectCategoryStatusTestSuite struct {
	suite.Suite
	handler                       command.ChangeProjectCategoryStatusHandler
	mockCtrl                      *gomock.Controller
	mockProjectCategoryRepository *mockproject.MockProjectCategoryRepository
	mockCachingRepository         *mockproject.MockCachingRepository
	mockService                   *mockproject.MockService
}

func (s *changeProjectCategoryStatusTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCategoryRepository = mockproject.NewMockProjectCategoryRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)
	s.handler = command.NewChangeProjectCategoryStatusHandler(s.mockProjectCategoryRepository, s.mockCachingRepository, s.mockService)
}

func (s *changeProjectCategoryStatusTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *changeProjectCategoryStatusTestSuite) Test_1_Success() {
	projectID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()
	status := domain.StatusInactive.String()

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	s.mockProjectCategoryRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCategoriesByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, performerID, projectID, categoryID, dto.ChangeProjectCategoryStatusRequest{
		Status: status,
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectCategoryStatusTestSuite) Test_1_Success_StatusNotChanged() {
	projectID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Status: domain.StatusInactive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, performerID, projectID, categoryID, dto.ChangeProjectCategoryStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectCategoryStatusTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.ProjectNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, uuid.New(), "invalid-id", uuid.New(), dto.ChangeProjectCategoryStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *changeProjectCategoryStatusTestSuite) Test_2_Fail_InvalidCategoryID() {
	projectID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCategory)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, performerID, projectID, "invalid-category-id", dto.ChangeProjectCategoryStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.InvalidCategory, err)
}

func (s *changeProjectCategoryStatusTestSuite) Test_2_Fail_NotOwner() {
	projectID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Common.NotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, performerID, projectID, categoryID, dto.ChangeProjectCategoryStatusRequest{
		Status: domain.StatusInactive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.NotFound, err)
}

func (s *changeProjectCategoryStatusTestSuite) Test_2_Fail_InvalidStatus() {
	projectID := uuid.New()
	categoryID := uuid.New()
	performerID := uuid.New()

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectCategoryStatus(ctx, performerID, projectID, categoryID, dto.ChangeProjectCategoryStatusRequest{
		Status: "invalid-status",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.InvalidStatus, err)
}

func TestChangeProjectCategoryStatusTestSuite(t *testing.T) {
	suite.Run(t, new(changeProjectCategoryStatusTestSuite))
}
