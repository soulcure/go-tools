package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   1024,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}

			//验证密码
			/*if _, authErr := conn.Do("AUTH", "123456"); authErr != nil {
				return nil, authErr
			}*/

			return conn, err
		},
	}
}

func Set(key, value string) {
	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	_, err := c.Do("Set", key, value)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func Get(key string) string {
	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	r, err := redis.String(c.Do("Get", key))
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return r
}
