package controllers

import (
	"GoStore/internal/config"
	"GoStore/pkg/store"
)

type Controller struct {
	store             *store.Store
	JWTKey            []byte
	userController    *UserController
	productController *ProductController
	orderController   *OrderController
	sessionController *SessionController
	chatController    *ChatController
}

func (controller *Controller) UserController() *UserController {
	if controller.userController != nil {
		return controller.userController
	}

	controller.userController = &UserController{controller: controller}

	return controller.userController
}

func (controller *Controller) ProductController() *ProductController {
	if controller.productController != nil {
		return controller.productController
	}

	controller.productController = &ProductController{controller: controller}

	return controller.productController
}

func (controller *Controller) OrderController() *OrderController {
	if controller.orderController != nil {
		return controller.orderController
	}

	controller.orderController = &OrderController{controller: controller}

	return controller.orderController
}

func (controller *Controller) SessionController() *SessionController {
	if controller.sessionController != nil {
		return controller.sessionController
	}

	controller.sessionController = &SessionController{controller: controller}

	return controller.sessionController
}

func (controller *Controller) ChatController() *ChatController {
	if controller.chatController != nil {
		return controller.chatController
	}

	controller.chatController = &ChatController{controller: controller}

	return controller.chatController
}

func New(cg *config.Config) (*Controller, error) {
	s, err := store.Open(cg)
	if err != nil {
		return nil, err
	}

	return &Controller{store: s, JWTKey: cg.JWTKey}, nil
}
