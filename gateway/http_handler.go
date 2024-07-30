package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"github.com/wesleybruno/golang-grpc-micro-service/gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(client gateway.OrdersGateway) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {

	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.handleGetOrder)

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

	o, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
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

	res := &CreateOrderRequest{
		Order:         o,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s", o.CustomerID, o.ID),
	}

	common.WriteJSON(w, http.StatusOK, res)

}

func (h *handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")
	orderId := r.PathValue("orderID")

	o, err := h.gateway.GetOrderById(r.Context(), customerId, orderId)

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
