package service

import (
	"context"
	"errors"
	"go-study/day01/internal/domain"
	"go-study/day01/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// 错误
var (
	ErrUserDulicateEmail     = repository.ErrUserDulicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

// 初始化
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// 注册
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

// 登录
func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	//先查是否有这个用户
	u, err := svc.repo.FindByEmail(ctx, user.Email)
	//没有这个用户
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		//返回账号或者密码错误
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
