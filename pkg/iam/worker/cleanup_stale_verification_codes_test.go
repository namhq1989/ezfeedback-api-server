package worker_test

import (
	"context"
	"testing"

	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type cleanupStaleVerificationCodesTestSuite struct {
	suite.Suite
	handler                        worker.CleanupStaleVerificationCodesHandler
	mockCtrl                       *gomock.Controller
	mockVerificationCodeRepository *mockiam.MockVerificationCodeRepository
}

func (s *cleanupStaleVerificationCodesTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *cleanupStaleVerificationCodesTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockVerificationCodeRepository = mockiam.NewMockVerificationCodeRepository(s.mockCtrl)

	s.handler = worker.NewCleanupStaleVerificationCodesHandler(s.mockVerificationCodeRepository)
}

func (s *cleanupStaleVerificationCodesTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *cleanupStaleVerificationCodesTestSuite) Test_1_Success() {
	// mock
	s.mockVerificationCodeRepository.EXPECT().
		CleanupStale(gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.CleanupStaleVerificationCodes(ctx, domain.QueueCleanupStaleVerificationCodesPayload{})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestCleanupStaleVerificationCodesTestSuite(t *testing.T) {
	suite.Run(t, new(cleanupStaleVerificationCodesTestSuite))
}
