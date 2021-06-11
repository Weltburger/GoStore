package controllers

import (
	"GoStore/pkg/models"
	"GoStore/pkg/statuserror"
	"context"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type OrderController struct {
	controller *Controller
}

func (orderController *OrderController) CreateOrder(c *gin.Context) {
	param := struct {
		ProductUUID string `json:"product_uuid" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
		Email       string `json:"email" binding:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	order, err := orderController.controller.store.OrderRepository().InsertOrder(ctx, &models.Order{
		ProductUUID: uuid.Must(uuid.FromString(param.ProductUUID)),
		Price:       param.Price,
		Quantity:    param.Quantity,
		Email:       param.Email,
		Status:      1,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "order has been placed",
		"value":   order.UUID,
	})
}

func (orderController *OrderController) AllOrders(c *gin.Context) {

	param := struct {
		Take int64 `form:"take,default=10"`
		Skip int64 `form:"skip,default=0"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(err)
	}

	ctx := context.Background()

	orders, total, err := orderController.controller.store.OrderRepository().AllOrders(ctx, param.Take, param.Skip)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"metadata": gin.H{
			"take":  param.Take,
			"skip":  param.Skip,
			"total": total,
		},
	})
}

func (orderController *OrderController) OrdersByEmail(c *gin.Context) {
	param := struct {
		Email string `json:"email" binding:"required"`
	}{}
	if err := c.ShouldBind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	orders, total, err := orderController.controller.store.OrderRepository().OrdersByEmail(ctx, &models.Order{
		Email: param.Email,
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"metadata": gin.H{
			"total": total,
		},
	})
}

func (orderController *OrderController) UpdateOrderStatus(c *gin.Context) {
	param := struct {
		UUID   uuid.UUID `json:"uuid" binding:"required"`
		Status int       `json:"status" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	err := orderController.controller.store.OrderRepository().UpdateOrderStatus(ctx, &models.Order{
		UUID:   param.UUID,
		Status: param.Status,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "order status has been updated",
	})
}

func (orderController *OrderController) UpdateOrder(c *gin.Context) {
	param := struct {
		UUID        uuid.UUID `json:"uuid" binding:"required"`
		ProductUUID uuid.UUID `json:"product_uuid" binding:"required"`
		Price       int       `json:"price" binding:"required"`
		Quantity    int       `json:"quantity" binding:"required"`
		Email       string    `json:"email" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	err := orderController.controller.store.OrderRepository().UpdateOrder(ctx, &models.Order{
		UUID:        param.UUID,
		ProductUUID: param.ProductUUID,
		Price:       param.Price,
		Quantity:    param.Quantity,
		Email:       param.Email,
		Status:      1,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "order has been updated",
	})
}

func (orderController *OrderController) DeleteOrder(c *gin.Context) {
	param := struct {
		UUID   uuid.UUID `json:"uuid" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	err := orderController.controller.store.OrderRepository().DeleteOrder(ctx, &models.Order{
		UUID:   param.UUID,
		Status: 0,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "order has been deleted",
	})
}
