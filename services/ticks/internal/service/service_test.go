package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/client"
	"github.com/stretchr/testify/suite"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	vdb       vdb.Port
	client    *client.GrpcClient
	closeTest func() error
}

func (suite *ServiceSuite) SetupTest() {
	defer tmpEnvVar("CRYPTELLATION_TICKS_GRPC_URL", ":9005")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_TICKS_GRPC_URL")
	grpcController := grpc.New(a)
	suite.NoError(grpcController.RunOnAddr(rpcUrl))

	ok := waitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.New()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		err = closeClient()
		go grpcController.Stop() // TODO: remove goroutine
		closeApplication()
		return err
	}

	vdb, err := redis.New()
	suite.Require().NoError(err)
	suite.vdb = vdb
}

func (suite *ServiceSuite) AfterTest(suiteName, testName string) {
	err := suite.closeTest()
	suite.NoError(err)
}

func (suite *ServiceSuite) TestListenSymbol() {
	ch, err := suite.client.Listen("SYMBOL")
	suite.Require().NoError(err)

	err = suite.client.Register(context.Background(), "mock_exchange", "SYMBOL")
	suite.Require().NoError(err)

	for i := int64(0); i < 50; i++ {
		t, ok := <-ch
		suite.Require().True(ok)
		suite.Require().Equal("mock_exchange", t.Exchange)
		suite.Require().Equal("SYMBOL", t.PairSymbol)
		suite.Require().Equal(float64(100+i), t.Price)
		suite.Require().WithinDuration(time.Unix(i, 0), t.Time, time.Microsecond)
	}
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}

func waitForPort(address string) bool {
	waitChan := make(chan struct{})

	go func() {
		for {
			conn, err := net.DialTimeout("tcp", address, time.Second)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if conn != nil {
				waitChan <- struct{}{}
				return
			}
		}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-waitChan:
		return true
	case <-timeout:
		return false
	}
}
