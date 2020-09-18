package controller

import (
	"distributed/sms/src/model"
	"distributed/sms/src/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

type TemporaryDataStorage struct {
	distributedCache *utils.DistributedCache
}

func NewTemporaryDataStorage(distributedCache *utils.DistributedCache) *TemporaryDataStorage {
	return &TemporaryDataStorage{
		distributedCache: distributedCache,
	}
}

func (t *TemporaryDataStorage) SetUserCertification(userId string, certification *model.Certification) error {
	certificationKey := fmt.Sprintf("userCertification:%s", userId)
	if err := t.distributedCache.Set(certificationKey, certification); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (t *TemporaryDataStorage) GetUserCertification(userId string) (*model.Certification, error) {
	certificationKey := fmt.Sprintf("userCertification:%s", userId)
	reply, err := t.distributedCache.Get(certificationKey)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	replyBytes, err := redis.Bytes(reply, err)
	if replyBytes == nil {
		return nil, nil
	}
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	certification := &model.Certification{}
	if err := json.Unmarshal(replyBytes, certification); err != nil {
		logs.Error(err)
		return nil, err
	}
	return certification, nil
}

func (t *TemporaryDataStorage) Flush() error{
	if err := t.distributedCache.Flush();err!=nil{
		logs.Error(err)
		return err
	}
	return nil
}
