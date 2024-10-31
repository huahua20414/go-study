package repository

import (
	"context"
	"go-study/day01/internal/domain"
	"go-study/day01/internal/repository/cache"
	"go-study/day01/internal/repository/dao"
	"time"
)

var ErrUserDulicatePhone = dao.ErrUserDulicatePhone

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{dao: dao,
		cache: cache}
}

// 修改密码
func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	return r.dao.Updates(ctx, dao.User{Phone: user.Phone, Password: user.Password})
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Phone:    user.Phone,
		Password: user.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	user := dao.User{Phone: u.Phone, Password: u.Password}
	err := r.dao.Insert(ctx, &user)
	if err != nil {
		return err
	}
	u.Id = user.Id
	u.Ctime = user.Ctime
	u.Utime = user.Utime
	return r.cache.Set(ctx, u)
}
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		//有数据
		return u, nil
	}
	//缓存没有
	if err == cache.ErrUserNotFound {
		//去数据库里面查
		u, err := r.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		user := domain.User{
			Id:       u.Id,
			Phone:    u.Phone,
			Password: u.Password,
			Ctime:    u.Ctime,
			Utime:    u.Utime,
		}
		err = r.cache.Set(ctx, user)
		if err != nil {
			//打日志做监控
		}
		return user, nil

	}
	//没数据,考虑redis崩掉，要不要去数据库查，做限流
	return domain.User{}, err

}

// 设置验证码缓存
func (r *UserRepository) SetVerification(ctx context.Context, user domain.User) error {
	user.Utime = time.Now().Unix()
	err := r.cache.Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// 获取验证码缓存
func (r *UserRepository) GetVerification(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := r.cache.GetVerification(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
