package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

var (
	//redisPoolInitCtx  sync.Once
	redisPoolInstance *redis.Pool
)

// 获得redis连接池
func GetRedisPool(addr string, maxIdle, maxActive, connectTimeout, readTimeout, writeTimeout int) *redis.Pool {

	//redisPoolInitCtx.Do(func() {

		redisPoolInstance = &redis.Pool{
			MaxIdle:   maxIdle,
			MaxActive: maxActive,
			Dial: func() (redis.Conn, error) {

				//c, err := redis.Dial("tcp", addr,
				//	redis.DialConnectTimeout(time.Duration(connectTimeout)*time.Millisecond),
				//	redis.DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
				//	redis.DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond))
				c, err := redis.Dial("tcp", addr)
				if err == nil && c != nil {
					return c, nil
				}

				return nil, errors.New("redisPool: cannot connect to any redis server")
			},
		}
	//})
	return redisPoolInstance
}
