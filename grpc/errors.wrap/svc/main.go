package main

import (
	"context"
	"flag"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jwenz723/errhandling/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func main() {
	svcName := "errhandling"

	fs := flag.NewFlagSet(svcName, flag.ExitOnError)
	grpcAddr := fs.String("grpc-addr", ":8082", "gRPC listen address")
	fs.Parse(os.Args[1:])

	logger, _ := zap.NewProduction()

	// Setup the server
	logger.Info("starting grpcSvc listener",
		zap.String("addr", *grpcAddr))
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Error("failed to start grpcSvc listener", zap.Error(err))
	}

	grpcSvc := grpcServer{}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)
	pb.RegisterOrdersServer(grpcServer, &grpcSvc)
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	// Do a client request to the server
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	s := pb.NewOrdersClient(conn)
	_, _ = s.NewOrder(context.TODO(), &pb.NewOrderRequest{CustomerID: "123"})

	lis.Close()
	time.Sleep(1 * time.Second)
}
