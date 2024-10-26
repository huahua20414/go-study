package repository

import (
	"context"
	"go-study/day01/internal/domain"
	"go-study/day01/internal/repository/dao"
)

var ErrUserDulicateEmail = dao.ErrUserDulicateEmail

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

// 修改密码
func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	return r.dao.Updates(ctx, dao.User{Email: user.Email, Password: user.Password})
}

func (r *UserRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	user, err := r.dao.FindByEmail(c, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *UserRepository) Create(c context.Context, u domain.User) error {
	return r.dao.Insert(c, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
