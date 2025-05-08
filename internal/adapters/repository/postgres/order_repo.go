package postgres

import (
	"time"

	"github.com/lioarce01/real-time-order-processing-system/internal/core/domain"
	"github.com/lioarce01/real-time-order-processing-system/internal/ports"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresOrderRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewOrderRepository(dsn string, logger *zap.Logger) (ports.OrderRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	// auto migrate the schema
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		logger.Error("AutoMigrate failed", zap.Error(err))
		return nil, err
	}

	//configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get database connection", zap.Error(err))
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &PostgresOrderRepository{db: db, logger: logger}, nil
}

func (r *PostgresOrderRepository) Save(order *domain.Order) error {
	// use GORM to save the order
	if err := r.db.Create(order).Error; err != nil {
		r.logger.Error("failed to save order", zap.Error(err))
		return err
	}

	// Verificar si el ID se asignó correctamente
	r.logger.Info("Order saved successfully", zap.Uint("order_id", order.ID))
	return nil
}

func (r *PostgresOrderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.Preload("Items").First(&order, id).Error; err != nil {
		r.logger.Error("failed to get order by id", zap.Error(err))
		return nil, err
	}
	return &order, nil
}
