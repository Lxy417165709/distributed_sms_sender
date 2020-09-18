package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

func GetRedisConnPool(redisHost string) *redis.Pool {
	return &redis.Pool{
		MaxActive:200,
		Dial:func()(redis.Conn,error) {
			conn,err := redis.Dial("tcp", redisHost)
			if err!=nil{
				logs.Error(err)
				return nil,err
			}
			return conn,nil
		},
		IdleTimeout:time.Second,
		Wait:true,
	}
}
