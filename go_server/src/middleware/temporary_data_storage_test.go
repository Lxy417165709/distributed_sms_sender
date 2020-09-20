package middleware

import (
	"distributed/sms/src/model"
	"distributed/sms/src/utils"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"testing"
)



func TestTemporaryDataStorage(t *testing.T) {
	userId := "1"
	redisHosts := []string{
		"120.26.162.39:20000",
		"120.26.162.39:20001",
		"120.26.162.39:20002",
		"120.26.162.39:20003",
	}
	hashLoop := NewHashLoop(30, map[int64]*redis.Pool{
		0:  utils.GetRedisConnPool(redisHosts[0]),
		8:  utils.GetRedisConnPool(redisHosts[1]),
		16: utils.GetRedisConnPool(redisHosts[2]),
		24: utils.GetRedisConnPool(redisHosts[3]),
	})
	distributedCache := NewDistributedCache(hashLoop)
	temporaryDataStorage := NewTemporaryDataStorage(distributedCache)
	if err := temporaryDataStorage.SetUserCertification(userId,&model.Certification{
		RemainSendTimes:8,
	});err!=nil{
		logs.Error(err)
		return
	}
	certification,err := temporaryDataStorage.GetUserCertification(userId)
	if err!=nil{
		logs.Error(err)
		return
	}
	logs.Info(certification)
}
