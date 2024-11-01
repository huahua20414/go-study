//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(day01-mysql:3306)/webook",
	},
	Redis: RedisConfig{
		Addr: "day01-redis:6379",
	},
}
