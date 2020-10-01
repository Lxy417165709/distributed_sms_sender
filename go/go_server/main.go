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
	"os"
	"strconv"
	"time"
)

const port = 8080

func init() {
	InitLogs()
}

func InitLogs() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func InitSingleSmsSender(smsSenderNum int) {
	const (
		sendTimesPerTerm = 200
		term = 60 * time.Second
		smsSenderMqTopic = "sms_sender"
		smsSenderLockKey = "smsSenderLock"
		userLockKey = "userLock"
		smsSenderOrderNumKey = "orderNum"
		userDistributedMutexRedisHost = "120.26.162.39:21000"
		smsDistributedMutexRedisHost = "120.26.162.39:21000"
		orderNumHost = "120.26.162.39:21000"
		userDistributedMutexRetryTime = 100 * time.Millisecond
		smsSenderDistributedMutexRetryTime = 100 * time.Millisecond
	)
	var (
		distributedCacheRedisHosts = []string{
			"120.26.162.39:20000",
			"120.26.162.39:20001",
			"120.26.162.39:20002",
			"120.26.162.39:20003",
		}
		kafkaHosts = []string{
			"qfstudio.net:15000",
			//"120.26.162.39:15000",
			//"120.26.162.39:15001",
			//"120.26.162.39:15002",
		}
	)

	var (
		hashLoop = middleware.NewHashLoop(30, map[int64]*redis.Pool{
			0:  utils.GetRedisConnPool(distributedCacheRedisHosts[0]),
			8:  utils.GetRedisConnPool(distributedCacheRedisHosts[1]),
			16: utils.GetRedisConnPool(distributedCacheRedisHosts[2]),
			24: utils.GetRedisConnPool(distributedCacheRedisHosts[3]),
		})
		distributedCache = middleware.NewDistributedCache(hashLoop)
		orderNumDistributedGenerator = middleware.NewOrderNumDistributedGenerator(
			utils.GetRedisConnPool(orderNumHost),
			smsSenderOrderNumKey,
		)
	)




	controller.SingleSmsSender = device.NewSmsSender(
		userDistributedMutexRedisHost,
		userLockKey,
		userDistributedMutexRetryTime,
		sendTimesPerTerm,
		term,
		middleware.NewTemporaryDataStorage(distributedCache),
		middleware.NewMqSender(kafkaHosts),
		smsSenderMqTopic,
		smsSenderNum,
		orderNumDistributedGenerator,
		middleware.NewDistributedMutex(smsDistributedMutexRedisHost, smsSenderLockKey, smsSenderDistributedMutexRetryTime),
	)
	if err := controller.SingleSmsSender.Flush(); err != nil {
		logs.Error(err)
		return
	}
}

func main() {
	InitSingleSmsSender(readSmsSenderNumFromOsArgs())
	r := gin.Default()
	r.GET("/test", controller.Test)
	r.POST("/send_msg", controller.SendMessage)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logs.Error("Running go http server failed. :|")
		return
	}
}

func readSmsSenderNumFromOsArgs() int {
	const defaultSmsSenderNum = 0
	smsSenderNum := defaultSmsSenderNum
	if len(os.Args) >= 2 {
		num, err := strconv.Atoi(os.Args[1])
		if err != nil {
			logs.Error(err)
			return smsSenderNum
		}
		smsSenderNum = num
	}
	return smsSenderNum
}
