package redis

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"time"
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

			log.Print("redis init success")
			return conn, err
		},
	}
}

func SetString(key, value string) {
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

func GetString(key string) string {
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

func SetInt(key string, value int) {
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

func GetInt(key string) int {
	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	r, err := redis.Int(c.Do("Get", key))
	if err != nil {
		logrus.Error(err)
		return 0
	}
	return r
}

func SetBool(key string, value bool) {
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

func GetBool(key string) bool {
	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	r, err := redis.Bool(c.Do("Get", key))
	if err != nil {
		logrus.Error(err)
		return false
	}
	return r
}

func SetFloat32(key string, value float32) {
	SetFloat64(key, float64(value))
}

func GetFloat32(key string) float32 {
	return float32(GetFloat64(key))
}

func SetFloat64(key string, value float64) {
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

func GetFloat64(key string) float64 {
	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	r, err := redis.Float64(c.Do("Get", key))
	if err != nil {
		logrus.Error(err)
		return 0
	}
	return r
}

//获取订单号
func GetOrderNum() string {
	var orderNo string

	c := pool.Get()
	defer func() {
		if err := pool.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	if num, err := redis.Int64(c.Do("INCR", "orderKey")); err != nil {
		logs.Error("redis orderKey get value error:", err)
		rand.Seed(time.Now().Unix())
		orderNo = "E" + time.Now().Format("20060102150405") + string(rand.Intn(100))
	} else {
		numStr := fmt.Sprintf("%04d", num)
		logs.Debug("order INCR:", numStr)
		orderNo = time.Now().Format("20060102150405") + numStr
	}

	return orderNo
}
