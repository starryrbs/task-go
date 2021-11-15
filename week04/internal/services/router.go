package services

import (
	"context"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctx context.Context, s *UserService) *gin.Engine {
	engine := gin.Default()
	engine.GET("/user", func(c *gin.Context) {
		s.ListUser(ctx, c)
	})

	engine.POST("/user", func(c *gin.Context) {
		s.CreateUser(ctx, c)
	})
	return engine
}
