package services

import (
	"context"
	"github.com/gin-gonic/gin"
	api "github.com/task-go/week04/api/user/v1"
	"github.com/task-go/week04/internal/biz"
	"time"
)

func NewUserService(user *biz.UserUseCase) *UserService {
	return &UserService{user}
}

func (s *UserService) ListUser(ctx context.Context, c *gin.Context) {
	users, err := s.user.List(ctx)
	if err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}
	rep := make([]*api.ListUserReply, 0)
	for _, user := range users {
		rep = append(rep, &api.ListUserReply{
			Id:   user.ID,
			Name: user.Name,
		})
	}
	c.JSON(200, rep)
}

func (s UserService) CreateUser(ctx context.Context, c *gin.Context) {
	// TODO 先写死，后面再优化

	err := s.user.Create(ctx, &biz.User{
		ID:        1,
		Name:      "哈哈",
		Email:     "",
		CreatedAt: time.Now(),
	})

	if err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}

	c.JSON(201, gin.H{"message": "创建成功"})
}
