package ports

import "github.com/lioarce01/real-time-order-processing-system/internal/core/domain"

type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id uint) (*domain.Order, error)
}
