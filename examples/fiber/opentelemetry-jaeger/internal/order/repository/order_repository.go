package repository

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/order/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OrderRepository represents a repository for managing Orders
type OrderRepository struct {
	Orders map[string]model.Order
	tracer trace.Tracer
}

// AddOrder adds a new Order to the repository
func (r *OrderRepository) AddOrder(ctx context.Context, order model.Order, parentSpan trace.Span) {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderRepository_AddOrder", trace.WithAttributes(
			attribute.String("order_id", order.OrderID),
			attribute.String("order_date", order.OrderDate),
			attribute.String("username", order.Username),
		))
	defer span.End()
	r.Orders[order.OrderID] = order
}

// GetOrderByOrderID retrieves a Order by their orderId
func (r *OrderRepository) GetOrderByOrderID(ctx context.Context, orderId string, parentSpan trace.Span) (*model.Order, error) {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderRepository_GetOrderByOrderID",
		trace.WithAttributes(
			attribute.String("order_id", orderId),
		),
	)
	defer span.End()
	for _, Order := range r.Orders {
		if Order.OrderID == orderId {
			return &Order, nil
		}
	}
	return nil, fmt.Errorf("order not found")
}

// GetAllOrders retrieves all Orders from the repository
func (r *OrderRepository) GetAllOrders(ctx context.Context, parentSpan trace.Span) map[string]model.Order {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderRepository_GetAllOrders",
	)
	defer span.End()
	return r.Orders
}

// DeleteOrder deletes an Order from the repository
func (r *OrderRepository) DeleteOrder(ctx context.Context, orderId string, parentSpan trace.Span) error {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderRepository_DeleteOrder",
		trace.WithAttributes(attribute.String("order_id", orderId)),
	)
	defer span.End()
	_, ok := r.Orders[orderId]
	if !ok {
		return fmt.Errorf("order not found for order id: %s", orderId)
	}

	delete(r.Orders, orderId)
	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order model.Order, parentSpan trace.Span) error {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderRepository_UpdateOrder",
		trace.WithAttributes(
			attribute.String("order_id", order.OrderID),
			attribute.String("order_date", order.OrderDate),
			attribute.String("username", order.Username),
		),
	)
	defer span.End()
	_, ok := r.Orders[order.OrderID]
	if !ok {
		return fmt.Errorf("order not found for order id: %s", order.OrderID)
	}
	r.Orders[order.OrderID] = order
	return nil
}

// NewOrderRepository creates a new OrderRepository instance
func NewOrderRepository(tracer trace.Tracer) *OrderRepository {
	return &OrderRepository{
		Orders: make(map[string]model.Order),
		tracer: tracer,
	}
}
