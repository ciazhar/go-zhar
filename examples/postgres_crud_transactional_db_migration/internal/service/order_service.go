package service

import (
	"context"
	"errors"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
)

type OrderServiceInterface interface {
	PlaceOrder(ctx context.Context, customerID int, items []model.OrderItem) (int, error)
	ProcessPayment(ctx context.Context, orderID int, method string, amount float64) error
	ShipOrder(ctx context.Context, orderID int, trackingNumber, carrier string) error
	MarkOrderDelivered(ctx context.Context, orderID int) error
}

type OrderService struct {
	orderRepository    *repository.PgxOrderRepository
	productRepository  *repository.PgxProductRepository
	paymentRepository  *repository.PgxPaymentRepository
	shipmentRepository *repository.PgxShipmentRepository
}

func NewOrderService(
	orderRepository *repository.PgxOrderRepository,
	productRepository *repository.PgxProductRepository,
	paymentRepository *repository.PgxPaymentRepository,
	shipmentRepository *repository.PgxShipmentRepository,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepository,
		productRepository:  productRepository,
		paymentRepository:  paymentRepository,
		shipmentRepository: shipmentRepository,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, customerID int, items []model.OrderItem) (int, error) {

	tx, err := s.orderRepository.BeginTransaction(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Create order
	orderID, err := s.orderRepository.CreateOrder(ctx, tx, customerID, "Pending")
	if err != nil {
		return 0, err
	}

	var totalAmount float64
	for _, item := range items {

		product, err := s.productRepository.GetProductByID(ctx, tx, item.ProductID)
		if err != nil {
			return 0, err
		}

		if product.Stock < item.Quantity {
			return 0, errors.New("insufficient stock")
		}

		totalPrice := float64(item.Quantity) * item.Price
		err = s.orderRepository.AddOrderItem(ctx, tx, orderID, item.ProductID, item.Quantity, item.Price, totalPrice)
		if err != nil {
			return 0, err
		}
		totalAmount += totalPrice

		// Adjust stock
		err = s.productRepository.AdjustStock(ctx, tx, item.ProductID, -item.Quantity)
		if err != nil {
			return 0, err
		}
	}

	// Update order total
	err = s.orderRepository.UpdateOrderTotal(ctx, tx, orderID, totalAmount)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (s *OrderService) ProcessPayment(ctx context.Context, orderID int, method string, amount float64) error {

	tx, err := s.orderRepository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.paymentRepository.ProcessPayment(ctx, tx, orderID, method, amount, "Success")
	if err != nil {
		return err
	}

	// Update order status to Completed
	err = s.orderRepository.UpdateOrderStatus(ctx, tx, orderID, "Completed")
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) ShipOrder(ctx context.Context, orderID int, trackingNumber, carrier string) error {
	return s.shipmentRepository.CreateShipment(ctx, orderID, trackingNumber, carrier, "Shipped")
}

func (s *OrderService) MarkOrderDelivered(ctx context.Context, orderID int) error {
	return s.shipmentRepository.UpdateShipmentStatus(ctx, orderID, "Delivered")
}
