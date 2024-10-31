package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-study/day01/internal/domain"
	"time"
)

var ErrUserNotFound = redis.Nil

type UserCache struct {
	client       redis.Cmdable
	expiration   time.Duration
	verification time.Duration
}

func NewUserCache(client redis.Cmdable, expiration time.Duration, verification time.Duration) *UserCache {
	return &UserCache{client: client, expiration: expiration, verification: verification}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	//数据不存在
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// 获取验证码缓存
func (cache *UserCache) GetVerification(ctx context.Context, u domain.User) (domain.User, error) {
	key := cache.keyPhone(u)
	//数据不存在
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		//数据不存在返回空
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	if err != nil {
		return domain.User{}, err
	}
	//如果redis有就返回查出来的对象
	return user, nil
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//生成唯一的键因为id唯一
	//如果id是0说明设置验证码缓存
	if u.Id == 0 {
		key := cache.keyPhone(u)
		return cache.client.Set(ctx, key, val, cache.verification).Err()
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

// 生成键的方法
func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
func (cache *UserCache) keyPhone(u domain.User) string {
	return fmt.Sprintf("user:info:%s:%s", u.CodeType, u.Phone)
}
