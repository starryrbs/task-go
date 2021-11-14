package biz

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserRepo interface {
	ListUser(ctx context.Context) ([]*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, id int64, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{repo}
}

func (uc *UserUseCase) List(ctx context.Context) (us []*User, err error) {
	return uc.repo.ListUser(ctx)
}

func (uc *UserUseCase) Get(ctx context.Context, id int64) (us *User, err error) {
	return uc.repo.GetUser(ctx, id)
}

func (uc *UserUseCase) Create(ctx context.Context, user *User) error {
	return uc.repo.CreateUser(ctx, user)
}

func (uc *UserUseCase) Update(ctx context.Context, id int64, user *User) error {
	return uc.repo.UpdateUser(ctx, id, user)
}

func (uc *UserUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteUser(ctx, id)
}
