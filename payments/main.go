package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/wesleybruno/golang-grpc-micro-service/common"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery/consul"
	"github.com/wesleybruno/golang-grpc-micro-service/payment/gateway"
	inmemProcessor "github.com/wesleybruno/golang-grpc-micro-service/payment/processor/inmem"
	"google.golang.org/grpc"
)

var (
	serviceName          = "payment"
	grpcAddress          = common.EnvString("GRPC_ADDRESS", "localhost:2001")
	consulAddr           = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser             = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass             = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost             = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.EnvString("RABBITMQ_PORT", "5672")
	httpAddr             = common.EnvString("HTTP_ADDR", "localhost:8081")
	endpointStripeSecret = common.EnvString("ENDPOINT_STRIPE_SECRET", "")
	jaegerAddr           = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		log.Fatal("could set global tracer")
	}

	// Service Discovery
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

	// Broker Connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	immenProcessor := inmemProcessor.NewInmem()
	gateway := gateway.NewGRPCGateway(registry)
	svc := NewService(immenProcessor, gateway)
	telemetryService := NewTelemetryMiddleware(svc)

	amqpConsumer := NewConsumer(telemetryService)
	go amqpConsumer.Listen(ch)

	// Http Server
	mux := http.NewServeMux()
	handler := NewPaymentHttpHandler(ch)
	handler.registerRoutes(mux)

	go func() {
		log.Println("Starting http server at:", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("Error to start http server", err.Error())
		}
	}()

	// GRPC Server
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
	defer l.Close()

	log.Println("New GRPC Server start at:", grpcAddress)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Error to start grpc server", err.Error())
	}
}
