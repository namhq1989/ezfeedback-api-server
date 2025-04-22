package command_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type updateProjectTestSuite struct {
	suite.Suite
	handler                command.UpdateProjectHandler
	mockCtrl               *gomock.Controller
	mockProjectRepository  *mockproject.MockProjectRepository
	mockProjectSettingRepo *mockproject.MockProjectSettingRepository
	mockCachingRepository  *mockproject.MockCachingRepository
}

func (s *updateProjectTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectSettingRepo = mockproject.NewMockProjectSettingRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.handler = command.NewUpdateProjectHandler(s.mockProjectRepository, s.mockProjectSettingRepo, s.mockCachingRepository)
}

func (s *updateProjectTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *updateProjectTestSuite) Test_1_Success() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	s.mockProjectRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockProjectSettingRepo.EXPECT().
		FindByProjectID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectSetting{ID: uuid.New()}, nil)

	s.mockProjectSettingRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteProjectSettingByProjectID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProject(ctx, uuid.New(), uuid.New(), dto.UpdateProjectRequest{
		Title:                  "Updated title",
		Description:            "Updated description",
		IsFeedbackPublic:       true,
		AllowAnonymousFeedback: false,
		EnableVoting:           true,
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *updateProjectTestSuite) Test_2_Fail_InvalidProjectID() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidProjectID)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProject(ctx, uuid.New(), "invalid-id", dto.UpdateProjectRequest{
		Title:       "Test",
		Description: "desc",
	})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
}

func (s *updateProjectTestSuite) Test_2_Fail_InvalidTitle() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProject(ctx, uuid.New(), uuid.New(), dto.UpdateProjectRequest{Title: ""})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidTitle, err)
}

func (s *updateProjectTestSuite) Test_2_Fail_InvalidDescription() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	longDesc := ""
	for i := 0; i < 2100; i++ {
		longDesc += "a"
	}

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateProject(ctx, uuid.New(), uuid.New(), dto.UpdateProjectRequest{
		Title:       "Test",
		Description: longDesc,
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidDescription, err)
}

//
// END OF CASES
//

func TestUpdateProjectTestSuite(t *testing.T) {
	suite.Run(t, new(updateProjectTestSuite))
}
