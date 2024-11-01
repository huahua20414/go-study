package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-study/webook/internal/domain"
	"time"
)

const verification = time.Minute * 5
const maxAttempts = 3 // 最大尝试次数

type CodeCache interface {
	Get(ctx context.Context, u domain.User) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	RemoveCode(ctx context.Context, u domain.User) error
	key(u domain.User) string
}

type RedisCodeCache struct {
	client       redis.Cmdable
	verification time.Duration
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client, verification: verification}
}

// 获取验证码缓存
func (cache *RedisCodeCache) Get(ctx context.Context, u domain.User) (domain.User, error) {
	key := cache.key(u)
	//查看次数=3直接删除
	attemptsKey := fmt.Sprintf("%s:attempts", key) // 尝试次数的键
	attempts, err := cache.client.Get(ctx, attemptsKey).Int()
	if err != nil {
		return domain.User{}, err
	}
	//尝试次数大于最大次数
	if attempts >= maxAttempts {
		//删除验证码和尝试次数
		if err := cache.RemoveCode(ctx, u); err != nil {
			return domain.User{}, err
		}
		return domain.User{}, nil
	}
	//尝试次数加一
	err = cache.client.Set(ctx, attemptsKey, attempts+1, time.Minute).Err() // 尝试次数存储1分钟
	if err != nil {
		return domain.User{}, err // 处理错误
	}
	val, err := cache.client.Get(ctx, key).Bytes()
	//系统错误
	if err != nil {
		return domain.User{}, err
	}
	//数据存在
	var user domain.User
	err = json.Unmarshal(val, &user)
	if err != nil {
		return domain.User{}, err
	}
	//如果redis有就返回查出来的对象
	return user, nil
}

func (cache *RedisCodeCache) Set(ctx context.Context, u domain.User) error {
	//解析user对象
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//生成唯一的键因为id唯一
	//设置验证码次数缓存
	key := cache.key(u)
	attemptsKey := fmt.Sprintf("%s:attempts", key) // 尝试次数的键
	//设置次数缓存1分钟过期
	if err := cache.client.Set(ctx, attemptsKey, 0, time.Minute).Err(); err != nil {
		return err
	}
	//设置验证码缓存
	return cache.client.Set(ctx, key, val, cache.verification).Err()
}

// 删除验证码
func (cache *RedisCodeCache) RemoveCode(ctx context.Context, u domain.User) error {
	key := cache.key(u)
	//删除验证码
	if err := cache.client.Del(ctx, key).Err(); err != nil && err != redis.Nil {
		return err
	}
	//删除验证码次数
	attemptsKey := fmt.Sprintf("%s:attempts", key)
	if err := cache.client.Del(ctx, attemptsKey).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (cache *RedisCodeCache) key(u domain.User) string {
	return fmt.Sprintf("user:info:%s:%s", u.CodeType, u.Phone)
}
