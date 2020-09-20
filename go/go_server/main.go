package main

import (
	"distributed/sms/src/controller"
	"distributed/sms/src/device"
	"distributed/sms/src/middleware"
	"distributed/sms/src/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"time"
)

const port = 8080

func init() {
	InitLogs()
	InitSingleSmsSender()
}

func InitLogs() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func InitSingleSmsSender() {
	redisHosts := []string{
		"120.26.162.39:20000",
		"120.26.162.39:20001",
		"120.26.162.39:20002",
		"120.26.162.39:20003",
	}
	kafkaHosts := []string{
		"120.26.162.39:15000",
		//"120.26.162.39:15001",
		//"120.26.162.39:15002",
	}
	distributedMutexRedisHost := "120.26.162.39:21000"
	hashLoop := middleware.NewHashLoop(30, map[int64]*redis.Pool{
		0:  utils.GetRedisConnPool(redisHosts[0]),
		8:  utils.GetRedisConnPool(redisHosts[1]),
		16: utils.GetRedisConnPool(redisHosts[2]),
		24: utils.GetRedisConnPool(redisHosts[3]),
	})

	distributedCache := middleware.NewDistributedCache(hashLoop)
	controller.SingleSmsSender = device.NewSmsSender(
		distributedMutexRedisHost,
		20,
		60*time.Second,
		middleware.NewTemporaryDataStorage(distributedCache),
		middleware.NewMqSender(kafkaHosts),
		"sms_sender",
	)
	if err := controller.SingleSmsSender.Flush(); err != nil {
		logs.Error(err)
		return
	}
}


func main() {
	r := gin.Default()
	r.GET("/test", controller.Test)
	r.POST("/send_msg", controller.SendMessage)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logs.Error("Running go http server failed. :|")
		return
	}
}

