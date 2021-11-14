package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/task-go/week04/internal/biz"
)

func NewUserService(user *biz.UserUseCase) *UserService {
	return &UserService{user}
}

func (s *UserService) ListUser(ctx context.Context, c *gin.Context) {
	s.user.List(ctx)
}
