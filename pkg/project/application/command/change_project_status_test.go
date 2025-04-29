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

type changeProjectStatusTestSuite struct {
	suite.Suite
	handler               command.ChangeProjectStatusHandler
	mockCtrl              *gomock.Controller
	mockProjectRepository *mockproject.MockProjectRepository
	mockCachingRepository *mockproject.MockCachingRepository
}

func (s *changeProjectStatusTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.handler = command.NewChangeProjectStatusHandler(s.mockProjectRepository, s.mockCachingRepository)
}

func (s *changeProjectStatusTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *changeProjectStatusTestSuite) Test_1_Success() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusInactive}, nil)

	s.mockProjectRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		SetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectsByUserID(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectStatus(ctx, performerID, projectID, dto.ChangeProjectStatusRequest{
		Status: domain.StatusActive.String(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectStatusTestSuite) Test_1_Success_StatusNotChanged() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusActive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectStatus(ctx, performerID, projectID, dto.ChangeProjectStatusRequest{
		Status: domain.StatusActive.String(),
	})
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *changeProjectStatusTestSuite) Test_2_Fail_ProjectNotFound() {
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectStatus(ctx, uuid.New(), uuid.New(), dto.ChangeProjectStatusRequest{
		Status: domain.StatusActive.String(),
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *changeProjectStatusTestSuite) Test_2_Fail_NotOwner() {
	projectID := uuid.New()

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: uuid.New(), Status: domain.StatusInactive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectStatus(ctx, uuid.New(), projectID, dto.ChangeProjectStatusRequest{
		Status: domain.StatusActive.String(),
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.NotFound, err)
}

func (s *changeProjectStatusTestSuite) Test_2_Fail_InvalidStatus() {
	var (
		projectID   = uuid.New()
		performerID = uuid.New()
	)

	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: projectID, UserID: performerID, Status: domain.StatusInactive}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeProjectStatus(ctx, performerID, projectID, dto.ChangeProjectStatusRequest{
		Status: "invalid-status",
	})
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidStatus, err)
}

//
// END OF CASES
//

func TestChangeProjectStatusTestSuite(t *testing.T) {
	suite.Run(t, new(changeProjectStatusTestSuite))
}
