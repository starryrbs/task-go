//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/task-go/week04/internal/biz"
	"github.com/task-go/week04/internal/conf"
	"github.com/task-go/week04/internal/data"
	"github.com/task-go/week04/internal/services"
)

func initApp(ctx context.Context, conf *conf.Conf) (*gin.Engine, error) {
	panic(wire.Build(data.ProviderSet, services.ProviderSet, biz.ProviderSet))
}
