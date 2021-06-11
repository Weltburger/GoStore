package controllers

import (
	"GoStore/pkg/models"
	"GoStore/pkg/statuserror"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

const (
	pattern = `^\w+@\w+\.\w+$`
)

type UserController struct {
	controller *Controller
}

func (userController *UserController) Login(c *gin.Context) {
	param := struct {
		Login    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusCodeBadParams, err))
	}

	ctx := context.Background()
	user, err := userController.controller.store.UserRepository().UserByEmail(ctx, &models.User{
		Email: param.Login,
	})
	if err != nil {
		panic(err)
	}

	if !user.CheckPassword(param.Password) {
		panic(statuserror.New(404, statuserror.StatusCodeBadParams, errors.New("bad password")))
	}

	accessToken, err := models.CreateToken(32)
	if err != nil {
		panic(err)
	}

	refreshToken, err := models.CreateToken(32)
	if err != nil {
		panic(err)
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"token": accessToken,
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"token": refreshToken,
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 144).Unix(),
	})

	userController.controller.store.Redis.Set(ctx,
		fmt.Sprintf("access_token:%s", accessToken),
		user.UUID,
		time.Hour*72)

	userController.controller.store.Redis.Set(ctx,
		fmt.Sprintf("refresh_token:%s", refreshToken),
		user.UUID,
		time.Hour*144)

	accessToken, err = jwtAccessToken.SignedString(userController.controller.JWTKey)
	if err != nil {
		panic(err)
	}

	refreshToken, err = jwtRefreshToken.SignedString(userController.controller.JWTKey)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (userController *UserController) Registration(c *gin.Context) {
	param := struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusCodeBadParams, err))
	}

	if matched, err := regexp.Match(pattern, []byte(param.Email)); err != nil || !matched {
		panic(statuserror.New(404, statuserror.StatusCodeBadParams, err))
	}

	ctx := context.Background()
	password, err := bcrypt.GenerateFromPassword([]byte(param.Password), 14)
	if err != nil {
		panic(statuserror.New(500, statuserror.StatusCodeServerErr, err))
	}

	_, err = userController.controller.store.UserRepository().InsertUser(ctx, &models.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: string(password),
	})
	if err != nil {
		panic(statuserror.New(500, statuserror.StatusCodeServerErr, err))
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "user created",
	})
}

func (userController *UserController) OK(c *gin.Context) {
	value, ok := c.Get("user")
	if !ok {
		panic(statuserror.NullUser)
	}

	user := value.(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"messages": user.Name,
	})
}

