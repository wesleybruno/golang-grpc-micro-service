package main

import (
	"context"
	"net"
	"time"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery/consul"
	"go.uber.org/zap"
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
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
	}

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
				logger.Error("Failed to health check", zap.Error(err))
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
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	store := NewOrderStore()
	svc := NewOrderService(store)
	telemetryService := NewTelemetryMiddleware(svc)

	NewGrRpcHandler(grpcServer, telemetryService, ch)

	consumer := NewConsumer(svc)
	go consumer.Listen(ch)

	logger.Info("Starting HTTP server", zap.String("port", grpcAddress))

	if err := grpcServer.Serve(l); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
