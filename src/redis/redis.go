package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type ConfigRedis struct {
	Redis InfoRedis `yaml:"redis"`
}

//数据库账号配置
type InfoRedis struct {
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

var (
	pool        *redis.Pool
	redisConfig ConfigRedis
)

func init() {
	path := "./conf/db.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("db conf file does not exist")
	}

	data, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(data, &redisConfig); err != nil {
		log.Panic("db conf yaml Unmarshal error ")
	}

	address := getConnURL(&redisConfig.Redis)

	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   1024,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			if err != nil {
				log.Panic("redis init error")
				return nil, err
			}

			password := redisConfig.Redis.Password
			if password != "" {
				//验证密码
				if _, authErr := conn.Do("AUTH", password); authErr != nil {
					return nil, authErr
				}
			}

			log.Print("redis connect to ", address)
			return conn, err
		},
	}
}

func getConnURL(info *InfoRedis) (url string) {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}

func DelKey(key string) (interface{}, error) {
	c := pool.Get()

	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return c.Do("DEL", key)
}

func ExpireKey(key string, seconds int64) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return c.Do("EXPIRE", key, seconds)
}

func GetKeys(pattern string) ([]string, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return redis.Strings(c.Do("KEYS", pattern))
}

func KeysByteSlices(pattern string) ([][]byte, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return redis.ByteSlices(c.Do("KEYS", pattern))
}

func SetBytes(key string, value []byte) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)

}

func GetBytes(key string) ([]byte, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.Bytes(c.Do("Get", key))
}

func SetString(key, value string) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)
}

func GetString(key string) (string, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.String(c.Do("Get", key))
}

func SetInt(key string, value int) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)
}

func GetInt(key string) (int, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.Int(c.Do("Get", key))
}

func SetInt64(key string, value int64) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)
}

func GetInt64(key string) (int64, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.Int64(c.Do("Get", key))
}

func SetBool(key string, value bool) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)
}

func GetBool(key string) (bool, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.Bool(c.Do("Get", key))
}

func SetFloat32(key string, value float32) (interface{}, error) {
	return SetFloat64(key, float64(value))
}

func GetFloat32(key string) (float32, error) {
	v, err := GetFloat64(key)
	return float32(v), err
}

func SetFloat64(key string, value float64) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return c.Do("Set", key, value)
}

func GetFloat64(key string) (float64, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	return redis.Float64(c.Do("Get", key))
}

// fieldValue 必须设置 tag ,如： Title  string `redis:"title"`
func SetStruct(key string, fieldValue interface{}) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return c.Do("HMSET", redis.Args{}.Add(key).AddFlat(fieldValue)...)
}

func GetStruct(key string, fieldValue interface{}) error {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	r, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		logrus.Error(err)
		return err
	}

	if e := redis.ScanStruct(r, fieldValue); e != nil {
		logrus.Error(err)
		return e
	}

	return err
}

func SetHashMap(key string, fieldValue map[string]interface{}) (interface{}, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return c.Do("HMSET", redis.Args{}.Add(key).AddFlat(fieldValue)...)
}

func GetHashMapString(key string) (map[string]string, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return redis.StringMap(c.Do("HGETALL", key))
}

func GetHashMapInt(key string) (map[string]int, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return redis.IntMap(c.Do("HGETALL", key))
}

func GetHashMapInt64(key string) (map[string]int64, error) {
	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	return redis.Int64Map(c.Do("HGETALL", key))
}

//获取订单号
func GetOrderNum() string {
	var orderNo string

	c := pool.Get()
	defer func() {
		if err := c.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	if num, err := redis.Int64(c.Do("INCR", "orderKey")); err != nil {
		logrus.Error("redis orderKey get value error:", err)
		rand.Seed(time.Now().Unix())
		orderNo = "E" + time.Now().Format("20060102150405") + string(rand.Intn(100))
	} else {
		numStr := fmt.Sprintf("%04d", num)
		logrus.Debug("order INCR:", numStr)
		orderNo = time.Now().Format("20060102150405") + numStr
	}

	return orderNo
}
