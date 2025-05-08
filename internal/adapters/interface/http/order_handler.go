package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lioarce01/real-time-order-processing-system/internal/core/domain"
	"github.com/lioarce01/real-time-order-processing-system/internal/core/service"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderSvc *service.OrderService
	logger   *zap.Logger
}

func NewOrderHandler(orderSvc *service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderSvc: orderSvc,
		logger:   logger,
	}
}

// Helper function to handle JSON responses
func (h *OrderHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if errMsg != "" {
		h.logger.Error(errMsg)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

// Helper function to extract and validate ID from the request URL
func (h *OrderHandler) extractOrderID(r *http.Request) (uint, error) {
	orderID, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		h.logger.Error("invalid order ID", zap.Error(err))
		return 0, err
	}
	return uint(orderID), nil
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("panic recovered in CreateOrder", zap.Any("error", r))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	var request struct {
		CustomerID uint               `json:"customer_id"`
		Items      []domain.OrderItem `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeJSONResponse(w, http.StatusBadRequest, nil, "Invalid request")
		return
	}

	h.logger.Info("Processing order", zap.Uint("customer_id", request.CustomerID))

	order, err := h.orderSvc.CreateOrder(request.CustomerID, request.Items)
	if err != nil {
		h.writeJSONResponse(w, http.StatusInternalServerError, nil, "Failed to create order")
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, order, "")
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL
	orderID, err := h.extractOrderID(r)
	if err != nil {
		h.writeJSONResponse(w, http.StatusBadRequest, nil, "Invalid order ID")
		return
	}

	// Retrieve order by ID via service
	order, err := h.orderSvc.GetOrderByID(orderID)
	if err != nil {
		h.writeJSONResponse(w, http.StatusInternalServerError, nil, "Failed to get order")
		return
	}

	// Respond with the order data
	h.writeJSONResponse(w, http.StatusOK, order, "")
}
