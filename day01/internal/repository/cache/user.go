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
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable, expiration time.Duration) *UserCache {
	return &UserCache{client: client, expiration: expiration}
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
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//生成唯一的键因为id唯一
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
