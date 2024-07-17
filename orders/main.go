package main

import (
	"context"
	"log"
	"net"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	"google.golang.org/grpc"
)

var (
	grpcAddress = common.EnvString("GRPC_ADDRESS", "localhost:2000")
)

func main() {

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
	defer l.Close()

	store := NewOrderStore()
	svc := NewOrderService(store)

	NewGrRpcHandler(grpcServer, *svc)
	svc.CreateOrder(context.Background())
	log.Println("New GRPC Server start at:", grpcAddress)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
}
