package grpc_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	mockproject "github.com/namhq1989/ezfeedback-api-server/internal/mock/project"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type onFeedbackCreatedTestSuite struct {
	suite.Suite
	handler             grpc.OnFeedbackCreatedHandler
	mockCtrl            *gomock.Controller
	mockQueueRepository *mockproject.MockQueueRepository
}

func (s *onFeedbackCreatedTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *onFeedbackCreatedTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockQueueRepository = mockproject.NewMockQueueRepository(s.mockCtrl)

	s.handler = grpc.NewOnFeedbackCreatedHandler(s.mockQueueRepository)
}

func (s *onFeedbackCreatedTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *onFeedbackCreatedTestSuite) Test_1_Success() {
	// mock
	s.mockQueueRepository.EXPECT().
		OnFeedbackCreated(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.OnFeedbackCreated(ctx, &projectpb.OnFeedbackCreatedRequest{
		TraceId:    "trace-id",
		ProjectId:  uuid.New(),
		CampaignId: uuid.New(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestOnFeedbackCreatedTestSuite(t *testing.T) {
	suite.Run(t, new(onFeedbackCreatedTestSuite))
}
