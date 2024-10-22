package controller

import (
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/model"
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/service"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// OrderController handles order-related requests
type OrderController struct {
	OrderService *service.OrderService
	tracer       trace.Tracer
}

// NewOrderController creates a new OrderController
func NewOrderController(
	orderService *service.OrderService,
) *OrderController {
	return &OrderController{
		OrderService: orderService,
		tracer:       otel.Tracer("OrderController"),
	}
}

// AddOrder @Summary Add a new order
func (uc *OrderController) AddOrder(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"OrderController_AddOrder")
	defer span.End()
	var order model.Order
	if err := c.BodyParser(&order); err != nil {
		span.RecordError(err)
		return c.Status(400).JSON(response.Response{
			Message: "Invalid request",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	uc.OrderService.AddOrder(c.Context(), order, span)
	return c.Status(201).JSON(response.Response{
		Message: "Order created",
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// GetOrderByOrderID is the handler for getting an order by their order id
func (uc *OrderController) GetOrderByOrderID(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"OrderController_GetOrderByOrderID",
	)
	defer span.End()
	orderId := c.Params("order_id")
	order, err := uc.OrderService.GetOrderByOrderID(c.Context(), orderId, span)
	if err != nil {
		span.RecordError(err)
		return c.Status(404).JSON(response.Response{
			Message: "Order not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(response.Response{
		Message: "Order found",
		Data:    order,
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// GetAllOrders is the handler for getting all orders
func (uc *OrderController) GetAllOrders(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"OrderController_GetAllOrders",
	)
	defer span.End()
	orders := uc.OrderService.GetAllOrders(c.Context(), span)
	return c.Status(200).JSON(response.Response{
		Message: "Orders found",
		Data:    orders,
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// DeleteOrder is the handler for deleting an order
func (uc *OrderController) DeleteOrder(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"OrderController_DeleteOrder",
	)
	defer span.End()
	orderId := c.Params("order_id")
	err := uc.OrderService.DeleteOrder(c.Context(), orderId, span)
	if err != nil {
		return c.Status(500).JSON(response.Response{
			Message: "Order not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(response.Response{
		Message: "Order deleted",
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// UpdateOrder is the handler for updating an order's information
func (uc *OrderController) UpdateOrder(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"OrderController_UpdateOrder",
	)
	defer span.End()
	var order model.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(response.Response{
			Message: "Invalid request",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	err := uc.OrderService.UpdateOrder(c.Context(), order, span)
	if err != nil {
		return c.Status(500).JSON(response.Response{
			Message: "Order not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(response.Response{
		Message: "Order updated",
		TraceID: span.SpanContext().TraceID().String(),
	})
}
