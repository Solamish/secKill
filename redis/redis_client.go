/**
*	redis客户端
*	封装了get\set\Lrang\Rpush等一系列常用操作
 */

package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)

// 定义redisPool对象连接ip、超时时长等属性
type Client struct {
	pool         *redis.Pool
	addr         string
	connTimeout  int
	readTimeout  int
	writeTimeout int
	maxIdle      int
	maxActive    int
	expireSecond int
}

var (
	RedisClent = &redisInit{
		RedisProxyClient: nil,
		RedisLayerClient: nil,
	}
)

type redisInit struct {
	RedisProxyClient *Client
	RedisLayerClient *Client
}

// main函数调用
func InitRedis() {
	RedisClent.RedisProxyClient = GetRedisInstance()
	RedisClent.RedisLayerClient = GetRedisInstance()

	initRedisProcess()
}

// TODO 参数配置化
func GetRedisInstance() *Client {
	addr := "127.0.0.1:6379"
	connTimeout := 100
	readTimeout := 50
	writeTimeout := 50
	maxIdle := 500
	maxActive := 1000
	expireSecond := 7000
	rc := new(Client)
	rc.pool = GetRedisPool(addr, maxIdle, maxActive, connTimeout, readTimeout, writeTimeout)
	if rc.pool == nil {
		log.Panic("get redis pool failed")
		return nil
	}

	rc.addr = addr
	rc.connTimeout = connTimeout
	rc.readTimeout = readTimeout
	rc.writeTimeout = writeTimeout
	rc.maxIdle = maxIdle
	rc.maxActive = maxActive
	rc.expireSecond = expireSecond

	return rc
}

func (rc *Client) Exists(key string) (bool, error) {
	c := rc.pool.Get()
	defer c.Close()
	exists, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false, err // handle error return from c.Do or type conversion error.
	}
	return exists, err
}

func (rc *Client) Expire(key string, expiresecond int) error {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Int(c.Do("EXPIRE", key, expiresecond))
	if err != nil {
		return err
	}

	if reply == 1 {
		return nil
	} else {
		return errors.New("redisClient: unexpected reply of expire")
	}
}

func (rc *Client) Del(key string) error {
	c := rc.pool.Get()
	defer c.Close()

	_, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		return err
	}
	//	if reply == 1 {
	//		return nil
	//	} else {
	// reply为0时说明key不存在
	//		return errors.New("redisClient: unexpected reply of del")
	//	}
	return nil
}

func (rc *Client) LPush(key, value string) error {
	c := rc.pool.Get()
	defer c.Close()
	_, err := redis.Int(c.Do("lpush", key, value))
	if err != nil {
		return err
	}

	// add redis key expire time.
	// ignore if error of expire command.
	_ = rc.Expire(key, rc.expireSecond)

	return nil
}

func (rc *Client) RPush(key, value string) error {
	c := rc.pool.Get()
	defer c.Close()

	_, err := redis.Int(c.Do("rpush", key, value))
	if err != nil {
		return err
	}

	// add redis key expire time.
	// ignore if error of expire command.
	_ = rc.Expire(key, rc.expireSecond)

	return nil
}

func (rc *Client) RPop(key string) (string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.String(c.Do("LPop", key))
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (rc *Client) BRPop(key string, timeOut int) ([]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Strings(c.Do("BRPop", key, timeOut))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (rc *Client) LRange(key string) ([]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Strings(c.Do("LRANGE", key, 0, -1))
	if err != nil {
		return nil, err
	}
	return reply, nil
}
