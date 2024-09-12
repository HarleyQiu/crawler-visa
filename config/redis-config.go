package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisConfig struct {
	Addr         string        `json:"addr"`          // 地址和端口
	Username     string        `json:"username"`      // 密码
	Password     string        `json:"password"`      // 密码
	DB           int           `json:"db"`            // 数据库索引
	DialTimeout  time.Duration `json:"dial_timeout"`  // 连接超时时间
	ReadTimeout  time.Duration `json:"read_timeout"`  // 读取超时时间
	WriteTimeout time.Duration `json:"write_timeout"` // 写入超时时间
	PoolSize     int           `json:"pool_size"`     // 连接池大小
}

func NewRedisClient(config *RedisConfig) *redis.Client {
	opts := &redis.Options{
		Addr:         config.Addr,
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		PoolTimeout:  30 * time.Second,
	}
	return redis.NewClient(opts)
}

func ConfigureRedis() *redis.Client {
	redisConfig := RedisConfig{
		Addr: "localhost:6379",
	}
	client := NewRedisClient(&redisConfig)
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
		return nil
	}
	log.Println("成功连接到Redis!")
	return client
}
