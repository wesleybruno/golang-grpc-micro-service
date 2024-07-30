package main

import pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}
