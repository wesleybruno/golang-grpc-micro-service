package main

import (
	"errors"
	"net/http"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)

}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {

	customerId := r.PathValue("customerId")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerId,
		Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {

		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, o)

}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("must be contains at least 1 item")
	}

	for _, i := range items {

		if i.ID == "" {
			return errors.New("must be contain ID")
		}

		if i.Quantity < 1 {
			return errors.New("must be contain Quantity")
		}

	}

	return nil

}
