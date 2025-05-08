package domain

import "time"

type OrderStatus string

const (
	OrderCreated    OrderStatus = "CREATED"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderCompleted  OrderStatus = "COMPLETED"
	OrderFailed     OrderStatus = "FAILED"
)

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	CustomerID uint        `gorm:"not null" json:"customer_id"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Total      float64     `json:"total"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
}
