package worker_test

import (
	"context"
	"testing"

	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type onFeedbackCreatedTestSuite struct {
	suite.Suite
	handler                       worker.OnFeedbackCreatedHandler
	mockCtrl                      *gomock.Controller
	mockProjectRepository         *mockproject.MockProjectRepository
	mockProjectCampaignRepository *mockproject.MockProjectCampaignRepository
	mockCachingRepository         *mockproject.MockCachingRepository
}

func (s *onFeedbackCreatedTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *onFeedbackCreatedTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectRepository = mockproject.NewMockProjectRepository(s.mockCtrl)
	s.mockProjectCampaignRepository = mockproject.NewMockProjectCampaignRepository(s.mockCtrl)
	s.mockCachingRepository = mockproject.NewMockCachingRepository(s.mockCtrl)

	s.handler = worker.NewOnFeedbackCreatedHandler(s.mockProjectRepository, s.mockProjectCampaignRepository, s.mockCachingRepository)
}

func (s *onFeedbackCreatedTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *onFeedbackCreatedTestSuite) Test_1_Success() {
	// mock
	s.mockProjectRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{
			ID: uuid.New(),
		}, nil)

	s.mockProjectRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockProjectCampaignRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.ProjectCampaign{
			ID: uuid.New(),
		}, nil)

	s.mockProjectCampaignRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		SetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		DeleteApiGetProjectByID(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.OnFeedbackCreated(ctx, domain.QueueOnFeedbackCreatedPayload{})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestOnFeedbackCreatedTestSuite(t *testing.T) {
	suite.Run(t, new(onFeedbackCreatedTestSuite))
}
