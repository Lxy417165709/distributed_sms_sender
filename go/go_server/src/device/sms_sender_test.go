package device

import (
	"distributed/sms/src/middleware"
	"distributed/sms/src/model"
	"distributed/sms/src/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)


func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func TestSingleSmsSender(t *testing.T) {
	smsSender := getSmsSender()
	if err := smsSender.Flush(); err != nil {
		logs.Error(err)
		return
	}
	userIds := make([]string, 0)
	for i := 0; i < 100; i++ {
		userIds = append(userIds, fmt.Sprintf("%d", i))
	}
	for {
		for _, userId := range userIds {
			go func(userId string) {
				result := smsSender.Send(&model.Message{
					SenderId: userId,
					ReceiverId: "-",
					Content: "test",
				})
				if result.IsSuccess {
					logs.Info("[%s] 号用户请求成功。返回消息：%s", userId, result.Msg)
				} else {
					//logs.Info("[%s] 号用户请求失败。返回消息：%s",userId,result.Msg)
				}
			}(userId)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestMulSmsSender(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	sendTimesPerTerm := 20
	smsSenders := make([]*SmsSender, 0)
	for i := 0; i < 10; i++ {
		smsSender := getSmsSender()
		if err := smsSender.Flush(); err != nil {
			logs.Error(err)
			return
		}
		smsSenders = append(smsSenders, smsSender)
	}
	var sendTimes int64 = 0
	userIds := make([]string, 0)
	for i := 0; i < 10; i++ {
		userIds = append(userIds, fmt.Sprintf("%d", i))
	}
	for {
		for _, userId := range userIds {
			go func(userId string) {
				senderNum := rand.Intn(len(smsSenders))
				result := smsSenders[senderNum].Send(&model.Message{
					SenderId: userId,
					ReceiverId: "-",
					Content: "test",
				})
				if result.IsSuccess {
					logs.Debug(atomic.AddInt64(&sendTimes, 1))
					logs.Info("[%d] 发送器 | [%s] 用户请求成功：%s", senderNum, userId, result.Msg)
					if sendTimes%int64(len(userIds)*sendTimesPerTerm) == 0 {
						var allTimes int64 = 0
						for i := 0; i < len(smsSenders); i++ {
							allTimes += smsSenders[i].GetHasSendTimes()
							logs.Info("[%d] 发送器 | 发送次数: %d", i, smsSenders[i].GetHasSendTimes())
						}
						logs.Debug("所有发送器发送的次数: %d", allTimes)
					}
				}

			}(userId)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func getSmsSender() *SmsSender{
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
	distributedMutexRedisHost := "127.0.0.1:6379"
	hashLoop := middleware.NewHashLoop(30, map[int64]*redis.Pool{
		0:  utils.GetRedisConnPool(redisHosts[0]),
		8:  utils.GetRedisConnPool(redisHosts[1]),
		16: utils.GetRedisConnPool(redisHosts[2]),
		24: utils.GetRedisConnPool(redisHosts[3]),
	})

	distributedCache := middleware.NewDistributedCache(hashLoop)
	return NewSmsSender(
		distributedMutexRedisHost,
		100,
		60*time.Second,
		middleware.NewTemporaryDataStorage(distributedCache),
		middleware.NewMqSender(kafkaHosts),
		"sms_sender",
	)
}
