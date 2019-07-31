package gredis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"gwyApp/pkg/setting"
	"time"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	logrus.Info("Setup is :", setting.RedisSetting.Host)
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// inct a key/value
func INCR(key string, time int64) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("INCR", key)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// sadd a key
func SAdd(key, value string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, value)
	if err != nil {
		return err
	}
	return nil
}

func SRem(key, value string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SREM", key, value)
	if err != nil {
		return err
	}
	return nil
}

func SIsmember(key, value string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("SISMEMBER", key, value))
}

func SPop(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("SPOP", key))

}

func ExpireAt(key string, expireMillSecond int64) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIREAT", key, expireMillSecond)
	if err != nil {
		return err
	}
	return nil
}
func LPush(key, value string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LPUSH", key, value))
}
func RPush(key string, value ...int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	//这里一个个插入
	//TODO:怎么用一条语句解决
	i := 0
	for _, v := range value {
		i++
		_, err := redis.Int(conn.Do("RPUSH", key, v))
		if err != nil {
			return 0, err
		}
	}
	return i, nil
}
func LIndex(key string, index int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LINDEX", key, index))
}
func LLen(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", key))
}
