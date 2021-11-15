package services

import (
	"github.com/google/wire"
	"github.com/task-go/week04/internal/biz"
)

type UserService struct {
	user *biz.UserUseCase
}

var ProviderSet = wire.NewSet(NewUserService, NewRouter)
