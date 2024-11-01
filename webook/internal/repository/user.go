package repository

import (
	"context"
	"go-study/webook/internal/domain"
	"go-study/webook/internal/repository/cache"
	"go-study/webook/internal/repository/dao"
	"time"
)

var ErrUserDulicatePhone = dao.ErrUserDulicatePhone

type UserRepository interface {
	Update(ctx context.Context, user domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	RemoveCode(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	SetVerification(ctx context.Context, user domain.User) error
	GetVerification(ctx context.Context, u domain.User) (domain.User, error)
}

type CachedUserRepository struct {
	dao       dao.UserDao
	cache     cache.UserCache
	codeCache cache.CodeCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache, codeCache cache.CodeCache) UserRepository {
	return &CachedUserRepository{dao: dao,
		cache: cache, codeCache: codeCache}
}

// 修改密码
func (r *CachedUserRepository) Update(ctx context.Context, user domain.User) error {
	return r.dao.Updates(ctx, dao.User{Phone: user.Phone, Password: user.Password})
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
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

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	user := dao.User{Phone: u.Phone, Password: u.Password}
	err := r.dao.Insert(ctx, &user)
	if err != nil {
		return err
	}
	u.Id = user.Id
	u.Ctime = user.Ctime
	u.Utime = user.Utime
	err = r.cache.Set(ctx, u)
	if err != nil {
		return err
	}
	//设置用户信息缓存成功,删除验证码
	u.CodeType = "register"
	return r.RemoveCode(ctx, u)

}

// user里要有codetype和phone的信息
func (r *CachedUserRepository) RemoveCode(ctx context.Context, u domain.User) error {
	if err := r.codeCache.RemoveCode(ctx, u); err != nil {
		return err
	}
	return nil
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
func (r *CachedUserRepository) SetVerification(ctx context.Context, user domain.User) error {
	user.Utime = time.Now().Unix()
	err := r.codeCache.Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// 获取验证码缓存
func (r *CachedUserRepository) GetVerification(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := r.codeCache.Get(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
