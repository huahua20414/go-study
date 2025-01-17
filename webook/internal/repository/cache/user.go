package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-study/webook/internal/domain"
	"time"
)

var ErrUserNotFound = redis.Nil

const expiration = time.Minute * 15

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	key(id int64) string
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{client: client, expiration: expiration}
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
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

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//生成唯一的键因为id唯一
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

// 生成键的方法
func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
