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

type getProjectByIDTestSuite struct {
	suite.Suite
	handler               query.GetProjectByIDHandler
	mockCtrl              *gomock.Controller
	mockCachingRepository *mockproject.MockCachingRepository
	mockService           *mockproject.MockService
}

func (s *getProjectByIDTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockproject.NewMockService(s.mockCtrl)

	s.handler = query.NewGetProjectByIDHandler(s.mockCachingRepository, s.mockService)
}

func (s *getProjectByIDTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectByIDTestSuite) Test_1_Success_FromCache() {
	cachedData := `{
		"projects": {
			"id": "d0388vpnlp83c24n98mg",
			"title": "EzFeedback",
			"slug": "ezfeedback-21hu1t",
			"status": "active",
			"description": "My first project description",
			"categories": [],
			"campaigns": [],
			"setting": {
			  "isFeedbackPublic": false,
			  "allowAnonymousFeedback": true,
			  "enableVoting": true
			},
			"stats": {
			  "totalFeedback": 10
			},
			"createdAt": "2025-04-22T00:46:39.916Z",
			"updatedAt": "2025-04-22T10:02:53.779Z"
		}
	}
	`

	s.mockCachingRepository.EXPECT().
		GetApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(&cachedData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetProjectByID(ctx, uuid.New(), uuid.New(), dto.GetProjectByIDRequest{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *getProjectByIDTestSuite) Test_1_Success_FromDB() {
	// mock data
	s.mockCachingRepository.EXPECT().
		GetApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockService.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	s.mockService.EXPECT().
		GetProjectSettingByProjectID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectSetting{ID: uuid.New()}, nil).
		AnyTimes()

	s.mockService.EXPECT().
		GetProjectCategoriesByProjectID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCategory{{ID: uuid.New()}}, nil).
		AnyTimes()

	s.mockService.EXPECT().
		GetProjectCampaignsByProjectID(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCampaign{{ID: uuid.New()}}, nil).
		AnyTimes()

	s.mockCachingRepository.EXPECT().
		SetApiGetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// call handler
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetProjectByID(ctx, uuid.New(), uuid.New(), dto.GetProjectByIDRequest{})

	// assertions
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetProjectByIDTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectByIDTestSuite))
}
