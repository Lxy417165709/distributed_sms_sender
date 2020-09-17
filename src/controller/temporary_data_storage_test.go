package controller

import (
	"distributed/sms/src/model"
	"github.com/astaxie/beego/logs"
	"testing"
)



func TestTemporaryDataStorage(t *testing.T) {
	userId := "1"
	temporaryDataStorage := NewTemporaryDataStorage(redisHost)
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
