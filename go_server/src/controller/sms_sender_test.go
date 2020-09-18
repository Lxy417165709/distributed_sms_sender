package controller

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)



func TestSingleSmsSender(t *testing.T) {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	smsSender := NewSmsSender(
		redisHost,
		100,
		time.Duration(1*time.Second),
		NewTemporaryDataStorage(redisHost),
	)
	if err := smsSender.Flush();err!=nil{
		logs.Error(err)
		return
	}
	userIds := make([]string,0)
	for i:=0;i<100;i++{
		userIds = append(userIds,fmt.Sprintf("%d",i))
	}

	for {
		for _,userId := range userIds{
			go func(userId string) {
				result := smsSender.Send(userId)
				if result.IsSuccess{
					logs.Info("[%s] 号用户请求成功。返回消息：%s",userId,result.Msg)
				}else{
					//logs.Info("[%s] 号用户请求失败。返回消息：%s",userId,result.Msg)
				}
			}(userId)
			time.Sleep(10*time.Millisecond)
		}
	}
	select{}
}


func TestMulSmsSender(t *testing.T) {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)

	sendTimesPerTerm := 10
	smsSenders := make([]*SmsSender,0)
	for i:=0;i<10;i++{
		smsSenders = append(smsSenders,NewSmsSender(
			redisHost,
			sendTimesPerTerm,
			time.Duration(10*time.Second),
			NewTemporaryDataStorage(redisHost),
		))
	}
	for i:=0;i<len(smsSenders);i++{
		if err := smsSenders[i].Flush();err!=nil{
			logs.Error(err)
			return
		}
	}
	var sendTimes int64 = 0
	userIds := make([]string,0)
	for i:=0;i<100;i++{
		userIds = append(userIds,fmt.Sprintf("%d",i))
	}

	rand.Seed(time.Now().UnixNano())
	for {
		for _,userId := range userIds{
			go func(userId string) {
				senderNum := rand.Intn(len(smsSenders))
				result := smsSenders[senderNum].Send(userId)
				if result.IsSuccess{
					logs.Debug(atomic.AddInt64(&sendTimes,1))
					logs.Info("[%d] 发送器 | [%s] 用户请求成功：%s",senderNum,userId,result.Msg)
					if sendTimes%int64(len(userIds)*sendTimesPerTerm)==0 {
						var allTimes int64 = 0
						for i:=0;i<len(smsSenders);i++{
							allTimes+=smsSenders[i].GetHasSendTimes()
							logs.Info("[%d] 发送器 | 发送次数: %d",i,smsSenders[i].GetHasSendTimes())
						}
						logs.Debug("所有发送器发送的次数: %d",allTimes)
					}
				}

			}(userId)
			time.Sleep(10*time.Millisecond)
		}
	}
	select{}
}

