package server

import (
	"GoStore/internal/config"
	"GoStore/internal/controllers"
	"GoStore/pkg/statuserror"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	sync.Once
}

func routing(conf *config.Config) (*gin.Engine, error) {
	controller, err := controllers.New(conf)
	if err != nil {
		return nil, err
	}

	router := gin.New()

	recovery := router.Group("/", func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case statuserror.IStatusError:
					c.JSON(e.HttpCode(), gin.H{
						"code":    e.StatusCode(),
						"message": e.Error(),
					})
				case error:
					c.JSON(500, gin.H{
						"error": e.Error(),
					})
				case string:
					c.JSON(500, gin.H{
						"error": e,
					})
				default:
					c.JSON(500, gin.H{
						"error": "undefined error",
					})
				}
				c.Abort()
			}
		}()
		c.Next()
	})

	recovery.POST("/login", controller.UserController().Login)
	recovery.POST("/registration", controller.UserController().Registration)
	recovery.POST("/select/product", controller.ProductController().SelectProduct)
	recovery.GET("/products", controller.ProductController().GetAll)

	recovery.GET("/ws", controller.ChatController().ListenSocket)
	//recovery.GET("/connect/to-chat", controller.ChatController().ConnectToChat)
	//recovery.GET("/select/chat", controller.ChatController().SelectChat)

	recovery.POST("/create/order", controller.OrderController().CreateOrder)

	sessionGroup := recovery.Group("/", controller.SessionController().CheckAuthorisation)

	sessionGroup.POST("/create/product", controller.ProductController().CreateProduct)
	sessionGroup.POST("/update/product", controller.ProductController().UpdateProduct)
	sessionGroup.POST("/delete/product", controller.ProductController().DeleteProduct)

	sessionGroup.GET("/orders", controller.OrderController().AllOrders)
	sessionGroup.POST("/orders-by-email", controller.OrderController().OrdersByEmail)
	sessionGroup.POST("/update/order", controller.OrderController().UpdateOrder)
	sessionGroup.POST("/update/order-status", controller.OrderController().UpdateOrderStatus)
	sessionGroup.POST("/delete/order", controller.OrderController().DeleteOrder)

	sessionGroup.GET("/show/chats", controller.ChatController().ShowChats)
	sessionGroup.GET("/select/chat", controller.ChatController().SelectChat)
	sessionGroup.POST("/delete/chat", controller.ChatController().RemoveHub)
	sessionGroup.POST("/delete/hubs", controller.ChatController().RemoveInactiveHubs)

	return router, nil
}

func New(path string) (*Server, error) {
	conf, err := config.New(path)
	if err != nil {
		return nil, err
	}

	router, err := routing(conf)
	if err != nil {
		return nil, err
	}

	return &Server{config: conf, router: router}, nil
}

func (server *Server) Start() error {
	serv := &http.Server{
		Addr:    server.config.Addr,
		Handler: server.router,
	}
	serv.SetKeepAlivesEnabled(true)
	ch := make(chan error, 1)
	go func() {
		if err := serv.ListenAndServe(); err != nil {
			ch <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case err := <-ch:
		return err
	case <-interrupt:
	}

	timeout, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer CancelFunc()

	if err := serv.Shutdown(timeout); err != nil {
		return err
	}

	return nil
}
