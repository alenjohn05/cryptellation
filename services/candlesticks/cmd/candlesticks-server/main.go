package main

import (
	"fmt"
	"os"

	"github.com/digital-feather/cryptellation/internal/genproto/candlesticks"
	"github.com/digital-feather/cryptellation/internal/server"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/service"
	"google.golang.org/grpc"
)

func run() int {
	application, err := service.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("creating application: %w", err))
		return 255
	}

	err = server.RunGRPCServer(func(server *grpc.Server) {
		svc := controllers.NewGrpcController(application)
		candlesticks.RegisterCandlesticksServiceServer(server, svc)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running application: %w", err))
		return 255
	}

	return 0
}

func main() {
	os.Exit(run())
}
