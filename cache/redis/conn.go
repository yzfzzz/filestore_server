package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
	redisHost = "127.0.0.1:6379"
)

// 创建redis连接池
func newRedisPool() *redis.Pool{
	return &redis.Pool{
		// 最大连接数
		MaxIdle: 50,
		// 可用的连接数
		MaxActive: 30,
		// 超时时间
		IdleTimeout: 300*time.Second,
		Dial: func()(redis.Conn,error){
			// 1.打开连接
			conn,err := redis.Dial("tcp",redisHost)
			if(err != nil){
				fmt.Println(err)
				return nil,err
			}
			// TODO: 访问认证
			
			return conn,nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if(time.Since(t) < time.Minute){
				return nil
			}
			_, err := c.Do("PING")
			return err
		},

	}
}

func init(){
	pool = newRedisPool()
}

func RedisPool() *redis.Pool{
	return pool
}