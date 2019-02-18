package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"test/extend/conf"
	"time"
)

var redisConn *redis.Pool

// GetRedisConn 获取 Redis 客户端连接
//想要操作redis就必须先获取操作redis的客户端对象
//理想状态下次reids对象应该提提供redis所有操作的方法，每个方法自己去维护池资源的释放以及获取，那么业务层的使用就会很舒服
// 我觉得所有的池都得这样玩。。。
func GetRedisConn() *redis.Pool {
	return redisConn
}

// Setup 创建 Redis 连接
func Setup() error {
	redisConn = &redis.Pool{
		MaxIdle: conf.RedisConf.MaxIdle,
		MaxActive: conf.RedisConf.MaxActive,
		IdleTimeout: conf.RedisConf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisConf.Host+":"+strconv.Itoa(conf.RedisConf.Port))
			if err != nil {
				return nil, err
			}
			// 验证密码
			if conf.RedisConf.Password != "" {
				if _, err := c.Do("AUTH", conf.RedisConf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			fmt.Println(conf.RedisConf.DBNum)
			// 选择数据库
			if _, err := c.Do("SELECT", conf.RedisConf.DBNum); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},

	}





	return nil
}

// Set 方法
func Set(key string, data string, seconds int) error {
	conn := GetRedisConn().Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return err
	}
	return nil
}

// Exists 方法
func Exists(key string) bool {
	conn := GetRedisConn().Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// Get 方法
func Get(key string) (string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err == redis.ErrNil {
		return "", nil
	}
	return reply, nil
}

// Del 方法
func Del(key string) (bool, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// DelLike 方法
func DelLike(key string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err := Del(key)
		if err != nil {
			return err
		}
	}
	return nil
}


// HGETALL 方法
func HGETALL(key string) (s []string, err error) {
	conn := GetRedisConn().Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("hgetall", key))
	if err != nil {
		return
	}

	for _,va :=range reply{
		if b,ok := va.([]byte); ok  {
			s = append(s,string(b))
		}
	}

	return
}


//hget
func HGET(key string,keya string) (s string, err error){
	conn := GetRedisConn().Get()
	defer conn.Close()

	s, err = redis.String(conn.Do("HGET", key,keya))
	if err != nil {
		return
	}
	return
}

//hget
func HMGET(key ...interface{}) (s []string, err error){
	conn := GetRedisConn().Get()
	defer conn.Close()

	hmget,err :=conn.Do("hmget", key...)
	s, err = redis.Strings(hmget,err)

	return

}

//sadd
func SADD(key string,val ...string)(err error){
	conn := GetRedisConn().Get()
	defer conn.Close()

	conn.Do("SADD", key, val)
	return
}