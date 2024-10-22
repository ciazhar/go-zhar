package service

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/model"
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OrderService struct {
	orderRepository    *repository.OrderRepository
	userHTTPRepository *repository.UserHTTPRepository
	tracer             trace.Tracer
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	userHTTPRepository *repository.UserHTTPRepository,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepo,
		userHTTPRepository: userHTTPRepository,
		tracer:             otel.Tracer("OrderService"),
	}
}

// AddOrder adds a new order
func (s *OrderService) AddOrder(ctx context.Context, order model.Order, parentSpan trace.Span) {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderService_AddOrder", trace.WithAttributes(
			attribute.String("order_id", order.OrderID),
			attribute.String("order_date", order.OrderDate),
			attribute.String("username", order.Username),
		))
	s.orderRepository.AddOrder(ctx, order, span)
}

// GetOrderByOrderID retrieves a order by order id
func (s *OrderService) GetOrderByOrderID(ctx context.Context, orderId string, parentSpan trace.Span) (*model.OrderExtended, error) {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderService_GetOrderByOrderID",
		trace.WithAttributes(
			attribute.String("order_id", orderId),
		),
	)
	defer span.End()

	order, err := s.orderRepository.GetOrderByOrderID(ctx, orderId, span)
	if err != nil {
		return nil, err
	}

	user, err := s.userHTTPRepository.GetUserByUsername(ctx, span, order.Username)
	if err != nil {
		return nil, err
	}

	resp := model.OrderExtended{
		OrderID:   order.OrderID,
		OrderDate: order.OrderDate,
		User:      user,
	}

	return &resp, nil
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders(ctx context.Context, parentSpan trace.Span) (res []model.OrderExtended) {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderService_GetAllOrders",
	)
	defer span.End()

	orders := s.orderRepository.GetAllOrders(ctx, span)

	for i := range orders {
		user, err := s.userHTTPRepository.GetUserByUsername(ctx, span, orders[i].Username)
		if err != nil {
			continue
		}
		order := model.OrderExtended{
			OrderID:   orders[i].OrderID,
			OrderDate: orders[i].OrderDate,
			User:      user,
		}
		res = append(res, order)
	}

	return
}

// DeleteOrder deletes a order by order id
func (s *OrderService) DeleteOrder(ctx context.Context, orderId string, parentSpan trace.Span) error {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderService_DeleteOrder",
		trace.WithAttributes(attribute.String("order_id", orderId)),
	)
	defer span.End()
	return s.orderRepository.DeleteOrder(ctx, orderId, span)
}

// UpdateOrder updates a order
func (s *OrderService) UpdateOrder(ctx context.Context, order model.Order, parentSpan trace.Span) error {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"OrderService_UpdateOrder",
		trace.WithAttributes(
			attribute.String("order_id", order.OrderID),
			attribute.String("order_date", order.OrderDate),
			attribute.String("username", order.Username),
		),
	)
	defer span.End()
	return s.orderRepository.UpdateOrder(ctx, order, span)
}
