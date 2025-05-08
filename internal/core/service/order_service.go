package service

import (
	"time"

	"github.com/lioarce01/real-time-order-processing-system/internal/core/domain"
	"github.com/lioarce01/real-time-order-processing-system/internal/ports"
	"go.uber.org/zap"
)

type OrderService struct {
	repo   ports.OrderRepository
	pub    ports.EventPublisher
	logger *zap.Logger
}

func NewOrderService(repo ports.OrderRepository, pub ports.EventPublisher, logger *zap.Logger) *OrderService {
	return &OrderService{repo: repo, pub: pub, logger: logger}
}

func (s *OrderService) CreateOrder(customerID uint, items []domain.OrderItem) (*domain.Order, error) {

	// Validate the order
	order := &domain.Order{
		CustomerID: customerID,
		Items:      items,
		Status:     domain.OrderCreated,
		CreatedAt:  time.Now(),
	}

	//compute total
	for _, item := range items {
		order.Total += item.Price * float64(item.Quantity)
	}

	// Persist order via repository (postgres)
	if err := s.repo.Save(order); err != nil {
		s.logger.Error("failed to save order",
			zap.Error(err),
			zap.Uint("customer_id", order.CustomerID),
			zap.Any("items", order.Items),
			zap.Float64("total", order.Total),
		)
		return nil, err
	}

	s.logger.Info("order saved", zap.Uint("order_id", order.ID))

	// Publish order created event
	if err := s.pub.Publish("orders", order); err != nil {
		// if publish fails, log warning
		s.logger.Warn("failed to publish order event", zap.Error(err))
	} else {
		s.logger.Info("order event published", zap.Uint("order_id", order.ID))
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(id uint) (*domain.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("failed to get order", zap.Error(err))
		return nil, err
	}
	return order, nil
}
