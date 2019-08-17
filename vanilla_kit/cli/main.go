package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/inContact/errhandling/vanilla_kit/kit"
	"google.golang.org/grpc"
	"os"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		panic (err)
	}
	svc := kit.NewGRPCClient(conn, logger)

	orderID, err := svc.NewOrder(context.TODO(), "123")

	logger.Log("orderID", orderID, "err", err)
}
