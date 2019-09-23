package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/inContact/orch-common/orchlog"
	orchlogflag "github.com/inContact/orch-common/orchlog/flag"
	"github.com/jwenz723/errhandling/pb"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"os"
	"path/filepath"
	"time"
)

// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func main() {
	svcName := "errhandling"

	cfg := struct {
		grpcAddr      string
		orchlogConfig orchlog.Config
	}{
		orchlogConfig: orchlog.Config{},
	}

	a := kingpin.New(filepath.Base(os.Args[0]), svcName)
	a.Flag("grpc-addr", "gRPC listen address.").Short('g').Default(":9884").StringVar(&cfg.grpcAddr)
	orchlogflag.AddFlags(a, &cfg.orchlogConfig)
	_, err := a.Parse(os.Args[1:])
	logger := orchlog.New(&cfg.orchlogConfig)

	var (
		endpointsLogger = log.With(logger,
			"component", "endpoint")
		gRPCLogger = log.With(logger,
			"component", "transport",
			"transport", "gRPC")
		gRPCClientLogger = log.With(logger,
			"component", "client",
			"transport", "gRPC")
	)

	svc := NewService()
	endpoints := NewSet(svc, endpointsLogger)
	grpcServer := NewGRPCServer(endpoints, gRPCLogger)

	// Setup the server
	grpcListener, err := net.Listen("tcp", cfg.grpcAddr)
	if err != nil {
		panic(err)
	}
	baseServer := grpc.NewServer()
	pb.RegisterOrdersServer(baseServer, grpcServer)
	go func() {
		_ = baseServer.Serve(grpcListener)
	}()

	// Do a client request to the server
	conn, err := grpc.Dial(cfg.grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	s := NewGRPCClient(conn, gRPCClientLogger)
	orderID, err := NewOrder(context.TODO(), "123")
	gRPCClientLogger.Log("orderID", orderID, "err", err)

	grpcListener.Close()
	time.Sleep(1 * time.Second)
}
