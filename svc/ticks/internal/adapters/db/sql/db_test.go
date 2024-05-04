package sql

import (
	"context"
	"os"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/ticks/deployments"
	"github.com/stretchr/testify/suite"
)

func TestSqlDatabaseSuite(t *testing.T) {
	suite.Run(t, new(SqlDatabaseSuite))
}

type SqlDatabaseSuite struct {
	suite.Suite
	db *Adapter
}

func (suite *SqlDatabaseSuite) SetupTest() {
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "ticks",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.db = db
}

func (suite *SqlDatabaseSuite) TestNewWithURIError() {
	defer tmpEnvVar("SQLDB_HOST", "")()

	var err error
	_, err = New(config.LoadSQL(nil))
	suite.Require().Error(err)
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
