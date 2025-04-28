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

type createProjectCategoryTestSuite struct {
	suite.Suite
	handler                       command.CreateProjectCategoryHandler
	mockCtrl                      *gomock.Controller
	mockProjectRepository         *mockproject.MockProjectRepository
	mockProjectCategoryRepository *mockproject.MockProjectCategoryRepository
	mockCachingRepository         *mockproject.MockCachingRepository
}

func (s *createProjectCategoryTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectCategoryRepository = mockproject.NewMockProjectCategoryRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.handler = command.NewCreateProjectCategoryHandler(s.mockProjectRepository, s.mockProjectCategoryRepository, s.mockCachingRepository)
}

func (s *createProjectCategoryTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createProjectCategoryTestSuite) Test_1_Success() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	s.mockProjectCategoryRepository.EXPECT().
		CountTotalByProjectID(gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	s.mockProjectCategoryRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectCategoriesByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCategory(ctx, performerID, projectID, dto.CreateProjectCategoryRequest{
		Name: "Category 1",
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.ID)
}

func (s *createProjectCategoryTestSuite) Test_2_Fail_ProjectNotFound() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCategory(ctx, uuid.New(), uuid.New(), dto.CreateProjectCategoryRequest{
		Name: "Category 1",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *createProjectCategoryTestSuite) Test_2_Fail_NotOwner() {
	projectID := uuid.New()
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: uuid.New(), Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCategory(ctx, uuid.New(), projectID, dto.CreateProjectCategoryRequest{
		Name: "Category 1",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.NotFound, err)
}

func (s *createProjectCategoryTestSuite) Test_2_Fail_CategoryLimitExceeded() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	s.mockProjectCategoryRepository.EXPECT().
		CountTotalByProjectID(gomock.Any(), gomock.Any()).
		Return(int64(100), nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCategory(ctx, performerID, projectID, dto.CreateProjectCategoryRequest{
		Name: "Category 1",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.CategoryLimitExceeded, err)
}

func (s *createProjectCategoryTestSuite) Test_5_Fail_InvalidName() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	s.mockProjectCategoryRepository.EXPECT().
		CountTotalByProjectID(gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProjectCategory(ctx, performerID, projectID, dto.CreateProjectCategoryRequest{
		Name: "a",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidName, err)
}

//
// END OF CASES
//

func TestCreateProjectCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(createProjectCategoryTestSuite))
}
