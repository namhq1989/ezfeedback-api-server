package command_test

import (
	"context"
	"testing"

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

type createProjectTestSuite struct {
	suite.Suite
	handler                      command.CreateProjectHandler
	mockCtrl                     *gomock.Controller
	mockProjectRepository        *mockproject.MockProjectRepository
	mockProjectSettingRepository *mockproject.MockProjectSettingRepository
	mockCachingRepository        *mockproject.MockCachingRepository
	mockBillingHub               *mockproject.MockBillingHub
}

func (s *createProjectTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectSettingRepository = mockproject.NewMockProjectSettingRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockBillingHub = mockproject.NewMockBillingHub(s.mockCtrl)
	s.handler = command.NewCreateProjectHandler(s.mockProjectRepository, s.mockProjectSettingRepository, s.mockCachingRepository, s.mockBillingHub)
}

func (s *createProjectTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createProjectTestSuite) Test_1_Success() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)

	s.mockProjectRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockProjectSettingRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectsByUserID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{
		Title:        "Test Project",
		Description:  "A project",
		Domain:       "test.com",
		PrimaryColor: "#FF0000",
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.ID)
}

func (s *createProjectTestSuite) Test_2_Fail_CanCreateProjectError() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(false, apperrors.Common.BadRequest)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{Title: "T"})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.BadRequest, err)
}

func (s *createProjectTestSuite) Test_2_Fail_ProjectLimitExceeded() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(false, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{Title: "T"})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Billing.ProjectLimitExceeded, err)
}

func (s *createProjectTestSuite) Test_2_Fail_NewProjectInvalidUserID() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, "invalid-id", dto.CreateProjectRequest{Title: "Test", Description: "desc"})
	assert.Nil(s.T(), resp)
	assert.NotNil(s.T(), err)
}

func (s *createProjectTestSuite) Test_2_Fail_NewProjectInvalidTitle() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{Title: ""})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidTitle, err)
}

func (s *createProjectTestSuite) Test_2_Fail_NewProjectInvalidDescription() {
	longDesc := ""
	for i := 0; i < 2100; i++ {
		longDesc += "a"
	}

	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{
		Title:       "Test",
		Description: longDesc,
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidDescription, err)
}

func (s *createProjectTestSuite) Test_2_Fail_NewProjectInvalidDomain() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)
	s.mockProjectRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{
		Title:  "Test",
		Domain: "invalid-domain",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidDomain, err)
}

func (s *createProjectTestSuite) Test_2_Fail_NewProjectInvalidPrimaryColor() {
	s.mockBillingHub.EXPECT().
		CanCreateProject(gomock.Any(), gomock.Any()).
		Return(true, nil)
	s.mockProjectRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateProject(ctx, uuid.New(), dto.CreateProjectRequest{
		Title:        "Test",
		PrimaryColor: "invalid-color",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidPrimaryColor, err)
}

//
// END OF CASES
//

func TestCreateProjectTestSuite(t *testing.T) {
	suite.Run(t, new(createProjectTestSuite))
}
