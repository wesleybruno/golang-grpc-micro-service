package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/wesleybruno/golang-grpc-micro-service/common"
	"github.com/wesleybruno/golang-grpc-micro-service/gateway/gateway"

	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery"
	"github.com/wesleybruno/golang-grpc-micro-service/common/discovery/consul"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		log.Fatal("could set global tracer")
	}
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceId, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err != nil {
				log.Fatalf("Error to sverify HealthCheck %s", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, serviceName)

	mux := http.NewServeMux()
	gateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(gateway)
	handler.registerRoutes(mux)

	log.Printf("Starting Server on PORT %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("Error to start server %s", err)
	}
}
