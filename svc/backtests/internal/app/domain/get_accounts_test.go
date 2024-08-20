package domain

import (
	"testing"

	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/events"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetAccountsSuite(t *testing.T) {
	suite.Run(t, new(GetAccountsSuite))
}

type GetAccountsSuite struct {
	suite.Suite
	operator     app.Backtests
	db           *db.MockPort
	Events       *events.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *GetAccountsSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.Events = events.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.Events, suite.candlesticks)
}

func (suite *GetAccountsSuite) TestHappyPass() {
	// TODO
}
