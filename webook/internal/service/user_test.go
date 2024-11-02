package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go-study/webook/internal/domain"
	"go-study/webook/internal/repository"
	repomocks "go-study/webook/internal/repository/mocks"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		wanterr  error
		wantuser domain.User
		user     domain.User
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).Return(domain.User{
					Phone:    "15930989101",
					Password: "$2a$10$zbe1e8D9iIw92y8/fVXnYefOYJd4jaW9/2zF2V8yTl2iYZ/3DbZbG",
				}, nil)
				return repo
			},
			wanterr: nil,
			user: domain.User{
				Phone:    "15930989101",
				Password: "123456",
			},
			wantuser: domain.User{
				Phone:    "15930989101",
				Password: "$2a$10$zbe1e8D9iIw92y8/fVXnYefOYJd4jaW9/2zF2V8yTl2iYZ/3DbZbG",
			},
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).Return(domain.User{
					Phone:    "15930989101",
					Password: "$2a$10$zbe1e8D9iIw92y8/fVXnYefOYJd4jaW9/2zF2V8yTl2iYZ/3DbZbG",
				}, nil)
				return repo
			},
			wanterr: errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
			user: domain.User{
				Phone:    "15930989101",
				Password: "12345",
			},
			wantuser: domain.User{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil, nil)
			u, err := svc.Login(context.Background(), tc.user)
			assert.Equal(t, tc.wanterr, err)
			assert.Equal(t, tc.wantuser, u)
		})
	}
}
func TestEn(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
