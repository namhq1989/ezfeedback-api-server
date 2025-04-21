package query_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type pingTestSuite struct {
	suite.Suite
	handler query.PingHandler
}

func (s *pingTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *pingTestSuite) setupApplication() {
	s.handler = query.NewPingHandler()
}

func (s *pingTestSuite) TearDownTest() {
}

//
// CASES
//

func (s *pingTestSuite) Test_1_Success() {
	// call
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.Ping(ctx, dto.PingRequest{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), true, resp.Success)
}

//
// END OF CASES
//

func TestPingTestSuite(t *testing.T) {
	suite.Run(t, new(pingTestSuite))
}
