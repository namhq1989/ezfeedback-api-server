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

type getProjectsTestSuite struct {
	suite.Suite
	handler               query.GetProjectsHandler
	mockCtrl              *gomock.Controller
	mockProjectRepository *mockproject.MockProjectRepository
	mockCachingRepository *mockproject.MockCachingRepository
	mockService           *mockproject.MockService
}

func (s *getProjectsTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = query.NewGetProjectsHandler(s.mockProjectRepository, s.mockCachingRepository, s.mockService)
}

func (s *getProjectsTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectsTestSuite) Test_1_Success_FromCache() {
	cachedData := `{
		"projects": [
			{
				"id": "d0388vpnlp83c24n98mg",
				"title": "EzFeedback",
				"slug": "ezfeedback-21hu1t",
				"status": "active",
				"stats": {
				  "totalFeedback": 10
				}
			}
		]
	}
	`

	s.mockCachingRepository.EXPECT().
		GetApiGetProjectsByUserID(gomock.Any(), gomock.Any()).
		Return(&cachedData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetProjects(ctx, uuid.New(), dto.GetProjectsRequest{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Projects))
}

func (s *getProjectsTestSuite) Test_2_Success_FromDB() {
	// mock data
	s.mockCachingRepository.EXPECT().
		GetApiGetProjectsByUserID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockProjectRepository.EXPECT().
		FindByUserID(gomock.Any(), gomock.Any()).
		Return([]domain.Project{{ID: uuid.New()}}, nil)

	s.mockService.EXPECT().
		GetProjectSettingByProjectID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectSetting{ID: uuid.New()}, nil).
		AnyTimes()

	s.mockService.EXPECT().
		GetProjectCategoriesByProjectID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCategory{{ID: uuid.New()}}, nil).
		AnyTimes()

	s.mockCachingRepository.EXPECT().
		SetApiGetProjectsByUserID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// call handler
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetProjects(ctx, uuid.New(), dto.GetProjectsRequest{})

	// assertions
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Projects))
}

//
// END OF CASES
//

func TestGetProjectsTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectsTestSuite))
}
