package data

import (
	"context"
	"github.com/task-go/week04/internal/biz"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{data: data}
}

func (up userRepo) ListUser(ctx context.Context) ([]*biz.User, error) {
	us, err := up.data.db.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.User, 0)
	for _, u := range us {
		rv = append(rv, &biz.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}
	return rv, nil
}

func (up *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	u, err := up.data.db.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (up userRepo) CreateUser(ctx context.Context, user *biz.User) error {
	_, err := up.data.db.User.
		Create().
		SetName(user.Name).
		SetCreatedAt(user.CreatedAt).
		SetID(user.ID).
		SetEmail(user.Email).
		Save(ctx)
	return err
}

func (up userRepo) UpdateUser(ctx context.Context, id int64, user *biz.User) error {
	u, err := up.data.db.User.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = u.Update().
		SetName(user.Name).
		Save(ctx)
	return err
}

func (up userRepo) DeleteUser(ctx context.Context, id int64) error {
	return up.data.db.User.DeleteOneID(id).Exec(ctx)
}
