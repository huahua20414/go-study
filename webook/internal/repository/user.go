package repository

import (
	"context"
	"go-study/webook/internal/domain"
	"go-study/webook/internal/repository/cache"
	"go-study/webook/internal/repository/dao"
	"time"
)

var ErrUserDulicatePhone = dao.ErrUserDulicatePhone

const maxAttempts = 3 // 最大尝试次数

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
	//注册成功删除验证码
	return r.RemoveCode(ctx, u)

}

// user里要有codetype和phone的信息
func (r *CachedUserRepository) RemoveCode(ctx context.Context, u domain.User) error {
	if err := r.codeCache.RemoveCode(ctx, cache.Code{
		CodeType: u.CodeType,
		Phone:    u.Phone,
	}); err != nil {
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
	err := r.codeCache.Set(ctx, cache.Code{
		CodeType:     user.CodeType,
		Phone:        user.Phone,
		Utime:        time.Now().Unix(),
		Verification: user.Verification,
	})
	if err != nil {
		return err
	}
	return nil
}

// 获取验证码缓存
func (r *CachedUserRepository) GetVerification(ctx context.Context, u domain.User) (domain.User, error) {
	code := cache.Code{
		CodeType: u.CodeType,
		Phone:    u.Phone,
		Utime:    0,
	}
	attempts, err := r.codeCache.GetAttempts(ctx, code)
	if err != nil {
		return domain.User{}, err
	}
	if attempts >= maxAttempts {
		//删除验证码和尝试次数
		if err := r.RemoveCode(ctx, u); err != nil {
			return domain.User{}, err
		}
		return domain.User{}, nil
	}
	//查看次数加一
	if err := r.codeCache.SetAttempts(ctx, code, attempts); err != nil {
		return domain.User{}, err
	}
	//获取验证码
	if code, err := r.codeCache.Get(ctx, code); err != nil {
		return domain.User{}, err
	} else {
		return domain.User{
			Utime:        code.Utime,
			Verification: code.Verification,
		}, nil
	}
}
