package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery/consul"
	"google.golang.org/grpc"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "orders"
	grpcAddress = common.EnvString("GRPC_ADDRESS", "localhost:2000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceId, serviceName, grpcAddress); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err != nil {
				log.Fatalf("Error to verify HealthCheck %s", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
	defer l.Close()

	store := NewOrderStore()
	svc := NewOrderService(store)

	NewGrRpcHandler(grpcServer, *svc, ch)
	svc.CreateOrder(context.Background())
	log.Println("New GRPC Server start at:", grpcAddress)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
}
