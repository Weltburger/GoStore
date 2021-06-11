package controllers

import (
	"GoStore/pkg/models"
	"GoStore/pkg/statuserror"
	"context"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type ProductController struct {
	controller *Controller
}

func (productController *ProductController) CreateProduct(c *gin.Context) {
	param := struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	value, ok := c.Get("user")
	if !ok {
		panic(statuserror.NullUser)
	}

	user := value.(*models.User)

	ctx := context.Background()

	product, err := productController.controller.store.ProductRepository().InsertProduct(ctx, &models.Product{
		User:        user,
		Name:        param.Name,
		Description: param.Description,
		Price:       param.Price,
		Quantity:    param.Quantity,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product has been created",
		"value":   product.UUID,
	})
}

func (productController *ProductController) DeleteProduct(c *gin.Context) {
	param := struct {
		UUID string `json:"uuid" binding:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	err := productController.controller.store.ProductRepository().DeleteProduct(ctx, &models.Product{
		UUID: uuid.Must(uuid.FromString(param.UUID)),
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product has been deleted",
	})
}

func (productController *ProductController) UpdateProduct(c *gin.Context) {
	param := struct {
		UUID        string `json:"uuid" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
	}{}
	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	err := productController.controller.store.ProductRepository().UpdateProduct(ctx, &models.Product{
		UUID:        uuid.Must(uuid.FromString(param.UUID)),
		Name:        param.Name,
		Description: param.Description,
		Price:       param.Price,
		Quantity:    param.Quantity,
	})

	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product has been updated",
	})
}

func (productController *ProductController) SelectProduct(c *gin.Context) {
	param := struct {
		UUID uuid.UUID `json:"uuid" binding:"required"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	ctx := context.Background()

	product, err := productController.controller.store.ProductRepository().ProductByUUID(ctx, &models.Product{
		UUID: param.UUID,//uuid.Must(uuid.FromString(param.UUID)),
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product has been selected",
		"value":   product,
	})
}

func (productController *ProductController) GetAll(c *gin.Context) {

	param := struct {
		Take int64 `form:"take,default=10"`
		Skip int64 `form:"skip,default=0"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(err)
	}

	ctx := context.Background()

	products, total, err := productController.controller.store.ProductRepository().AllProducts(ctx, param.Take, param.Skip)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"metadata": gin.H{
			"take": param.Take,
			"skip": param.Skip,
			"total": total,
		},
	})
}
