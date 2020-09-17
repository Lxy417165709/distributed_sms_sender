package controller

import (
	"distributed/sms/src/model"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

type TemporaryDataStorage struct{
	connPool *redis.Pool
}

func NewTemporaryDataStorage(redisHost string) *TemporaryDataStorage{
	pool := &redis.Pool{
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
	return &TemporaryDataStorage{
		connPool: pool,
	}
}

func (t *TemporaryDataStorage) SetUserCertification(userId string,certification *model.Certification)error {
	conn := t.connPool.Get()
	defer func() {
		if err := conn.Close();err!=nil{
			logs.Error(err)
		}
	}()
	certificationKey := fmt.Sprintf("userCertification:%s",userId)
	certificationBytes,err  := json.Marshal(certification)
	if err!=nil{
		logs.Error(err)
		return  err
	}
	if _,err := conn.Do("set",certificationKey,certificationBytes);err!=nil{
		logs.Error(err)
		return err
	}
	return nil
}
func (t *TemporaryDataStorage) GetUserCertification(userId string) (*model.Certification,error) {
	conn := t.connPool.Get()
	defer func() {
		if err := conn.Close();err!=nil{
			logs.Error(err)
		}
	}()
	certificationKey := fmt.Sprintf("userCertification:%s",userId)

	reply,err := conn.Do("get",certificationKey)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	replyBytes,err := redis.Bytes(reply,err)
	if replyBytes==nil{
		return nil,nil
	}
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	certification:=&model.Certification{}
	if err := json.Unmarshal(replyBytes,certification);err!=nil{
		logs.Error(err)
		return nil,err
	}
	return certification,nil
}
