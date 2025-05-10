package query_test

import (
	"context"
	"testing"

	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getProjectCollaboratorsTestSuite struct {
	suite.Suite
	handler     query.GetProjectCollaboratorsHandler
	mockCtrl    *gomock.Controller
	mockIAMHub  *mockproject.MockIAMHub
	mockService *mockproject.MockService
}

func (s *getProjectCollaboratorsTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockIAMHub = mockproject.NewMockIAMHub(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = query.NewGetProjectCollaboratorsHandler(s.mockIAMHub, s.mockService)
}

func (s *getProjectCollaboratorsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectCollaboratorsTestSuite) Test_1_Success() {
	// mock data
	s.mockService.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	s.mockService.EXPECT().
		GetProjectCollaboratorsByProjectID(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCollaborator{{ID: uuid.New()}}, nil)

	s.mockIAMHub.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil).
		AnyTimes()

	// call handler
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetProjectCollaborators(ctx, uuid.New(), uuid.New(), dto.GetProjectCollaboratorsRequest{})

	// assertions
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Collaborators))
}

//
// END OF CASES
//

func TestGetProjectCollaboratorsTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectCollaboratorsTestSuite))
}
