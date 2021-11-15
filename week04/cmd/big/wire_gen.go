// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/task-go/week04/internal/biz"
	"github.com/task-go/week04/internal/conf"
	"github.com/task-go/week04/internal/data"
	"github.com/task-go/week04/internal/services"
)

// Injectors from wire.go:

func initApp(ctx context.Context, conf2 *conf.Conf) (*gin.Engine, error) {
	dataData, err := data.NewData(ctx, conf2)
	if err != nil {
		return nil, err
	}
	userRepo := data.NewUserRepo(dataData)
	userUseCase := biz.NewUserUseCase(userRepo)
	userService := services.NewUserService(userUseCase)
	engine := services.NewRouter(ctx, userService)
	return engine, nil
}
