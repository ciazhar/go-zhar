package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, request model.OrderRequest) error
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetOrder(ctx context.Context, orderID int) (*model.Order, error)
	Delete(ctx context.Context, orderID int) error
}

type OrderService struct {
	orderRepository     repository.OrderRepository
	inventoryRepository repository.InventoryRepository
	paymentRepository   repository.PaymentRepository
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	inventoryRepository repository.InventoryRepository,
	paymentRepository repository.PaymentRepository,
) *OrderService {
	return &OrderService{
		orderRepository:     orderRepository,
		inventoryRepository: inventoryRepository,
		paymentRepository:   paymentRepository,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, request model.OrderRequest) error {

	tx, err := s.orderRepository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Step 1: Insert Order
	orderID, err := s.orderRepository.Create(ctx, tx, request.CustomerName)
	if err != nil {
		return err
	}

	// Step 2: Update Inventory
	for _, quantity := range request.Items {
		err = s.inventoryRepository.Update(ctx, tx, quantity.Quantity, quantity.ProductID)
		if err != nil {
			return err
		}
	}

	// Step 3: Process Payment
	err = s.paymentRepository.Create(ctx, tx, orderID, int(request.Amount), "Success")
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

func (s *OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	return s.orderRepository.GetAllOrders(ctx)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int) (*model.Order, error) {
	return s.orderRepository.GetOrderByID(ctx, orderID)
}

func (s *OrderService) Delete(ctx context.Context, orderID int) error {

	tx, err := s.orderRepository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Step 1: Delete Order
	err = s.orderRepository.DeleteOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Step 2: Delete Payment
	err = s.paymentRepository.Delete(ctx, tx, orderID)
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
