package common

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var DefaultRedisPool *redis.Pool

func InitRedis() {
	DefaultRedisPool = &redis.Pool{
		MaxIdle:     GetConfig().MustInt("redis", "maxIdle", 10),
		IdleTimeout: time.Duration(GetConfig().MustInt("redis", "idleTimeout", 10)) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", GetConfig().MustValue("redis", "host", "")+":"+GetConfig().MustValue("redis", "port", ""),
				redis.DialPassword(GetConfig().MustValue("redis", "password", "")),
				redis.DialConnectTimeout(time.Duration(GetConfig().MustInt("redis", "connectTimeout", 10))*time.Millisecond),
				redis.DialReadTimeout(time.Duration(GetConfig().MustInt("redis", "readTimeout", 10))*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(GetConfig().MustInt("redis", "writeTimeout", 10))*time.Millisecond))
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	RegisterComponent(&RedisPool{})
	DefLogger.Info("======finish redis init ")
}

type RedisPool struct {
}

func (r RedisPool) Destroy() {
	DefaultRedisPool.Close()
}
