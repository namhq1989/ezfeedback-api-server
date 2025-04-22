package command_test

import (
	"context"
	"testing"

	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type generateTokenTestSuite struct {
	suite.Suite
	handler           command.GenerateTokenHandler
	mockCtrl          *gomock.Controller
	mockJwtRepository *mockiam.MockJwtRepository
}

func (s *generateTokenTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *generateTokenTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockJwtRepository = mockiam.NewMockJwtRepository(s.mockCtrl)

	s.handler = command.NewGenerateTokenHandler(s.mockJwtRepository)
}

func (s *generateTokenTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *generateTokenTestSuite) Test_1_Success() {
	// mock data
	s.mockJwtRepository.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.Any()).
		Return("token", nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GenerateToken(ctx, dto.GenerateTokenRequest{
		UserID: "ref_id",
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.Token)
}

//
// END OF CASES
//

func TestGenerateTokenTestSuite(t *testing.T) {
	suite.Run(t, new(generateTokenTestSuite))
}
