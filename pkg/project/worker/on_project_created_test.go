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

type onProjectCreatedTestSuite struct {
	suite.Suite
	handler                           worker.OnProjectCreatedHandler
	mockCtrl                          *gomock.Controller
	mockProjectCollaboratorRepository *mockproject.MockProjectCollaboratorRepository
}

func (s *onProjectCreatedTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *onProjectCreatedTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectCollaboratorRepository = mockproject.NewMockProjectCollaboratorRepository(s.mockCtrl)

	s.handler = worker.NewOnProjectCreatedHandler(s.mockProjectCollaboratorRepository)
}

func (s *onProjectCreatedTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *onProjectCreatedTestSuite) Test_1_Success() {
	// mock
	s.mockProjectCollaboratorRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.OnProjectCreated(ctx, domain.QueueOnProjectCreatedPayload{
		ProjectID: uuid.New(),
		UserID:    uuid.New(),
	})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestOnProjectCreatedTestSuite(t *testing.T) {
	suite.Run(t, new(onProjectCreatedTestSuite))
}
