package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestHashLoop(t *testing.T) {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	redisHosts := []string{
		"120.26.162.39:20000",
		"120.26.162.39:20001",
		"120.26.162.39:20002",
		"120.26.162.39:20003",
	}
	hashLoop := NewHashLoop(30, map[int64]*redis.Pool{
		0:GetRedisConnPool(redisHosts[0]),
		8:GetRedisConnPool(redisHosts[1]),
		16:GetRedisConnPool(redisHosts[2]),
		24:GetRedisConnPool(redisHosts[3]),
	})
	nums := []int64{
		1,2,5,8,16,17,27,
	}
	for _,num := range nums{
		logs.Info(hashLoop.GetRedisPool(num))
	}
}
