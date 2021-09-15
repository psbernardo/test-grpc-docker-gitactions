package user

import (
	"context"
	"testing"

	"github.com/patrick/test-grpc-docker-gitactions/lib/dbtest"
	"github.com/patrick/test-grpc-docker-gitactions/proto/userpb"
	"github.com/stretchr/testify/suite"
)

type OutputsettlementTestSuite struct {
	dbtest.Suite
}

func TestOutputSettlement(t *testing.T) {
	suite.Run(t, new(OutputsettlementTestSuite))
}

func (s *OutputsettlementTestSuite) SetupSuite() {
	s.SetupDB()
}

func (s *OutputsettlementTestSuite) TearDownSuite() {
	s.TearDownDB()
}

func (s *OutputsettlementTestSuite) Test_Outputsettlement() {
	require := s.Require()
	ctx := context.Background()
	_, err := New().WithTx(&ctx).TestCall(&userpb.User{Name: "patrik test"})

	require.Nil(err)
}
