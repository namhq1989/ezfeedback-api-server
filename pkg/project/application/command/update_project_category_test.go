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

type updateProjectCategoryTestSuite struct {
	suite.Suite
	handler                       command.UpdateProjectCategoryHandler
	mockCtrl                      *gomock.Controller
	mockProjectCategoryRepository *mockproject.MockProjectCategoryRepository
	mockCachingRepository         *mockproject.MockCachingRepository
	mockService                   *mockproject.MockService
}

func (s *updateProjectCategoryTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCategoryRepository = mockproject.NewMockProjectCategoryRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)
	s.handler = command.NewUpdateProjectCategoryHandler(s.mockProjectCategoryRepository, s.mockCachingRepository, s.mockService)
}

func (s *updateProjectCategoryTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *updateProjectCategoryTestSuite) Test_1_Success() {
	var (
		projectID   = uuid.New()
		categoryID  = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Name: "OldName"}, nil)

	s.mockProjectCategoryRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCategoriesByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectsByUserID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCategory(ctx, performerID, projectID, categoryID, dto.UpdateProjectCategoryRequest{
		Name: "UpdatedName",
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *updateProjectCategoryTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.ProjectNotFound)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCategory(ctx, uuid.New(), "invalid-id", uuid.New(), dto.UpdateProjectCategoryRequest{
		Name: "Test",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *updateProjectCategoryTestSuite) Test_2_Fail_InvalidCategoryID() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCategory)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCategory(ctx, performerID, projectID, "invalid-category-id", dto.UpdateProjectCategoryRequest{
		Name: "Test",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Project.InvalidCategory, err)
}

func (s *updateProjectCategoryTestSuite) Test_2_Fail_InvalidName() {
	var (
		projectID   = uuid.New()
		categoryID  = uuid.New()
		performerID = uuid.New()
	)

	s.mockService.EXPECT().
		GetProjectCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCategory{ID: categoryID, ProjectID: projectID, Name: "OldName"}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProjectCategory(ctx, performerID, projectID, categoryID, dto.UpdateProjectCategoryRequest{
		Name: "",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.InvalidName, err)
}

func TestUpdateProjectCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(updateProjectCategoryTestSuite))
}
