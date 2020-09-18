package controller

import (
	"distributed/sms/src/model"
	"distributed/sms/src/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"sync/atomic"
	"time"
)

type SmsSender struct {
	redisHost                string
	userCertificationStorage *TemporaryDataStorage
	senderDistributeMutex    *utils.DistributedMutex
	userDistributedMutex     map[string]*utils.DistributedMutex
	sendTimesPerTerm         int
	term                     time.Duration
	hasSendTimes             int64
}

func NewSmsSender(redisHost string, sendTimesPerTerm int, term time.Duration, userCertificationStorage *TemporaryDataStorage) *SmsSender {
	return &SmsSender{
		redisHost:                redisHost,
		senderDistributeMutex:    utils.NewDistributedMutex(redisHost, "smsSenderLock", 1),
		sendTimesPerTerm:         sendTimesPerTerm,
		term:                     term,
		userCertificationStorage: userCertificationStorage,
		userDistributedMutex:     make(map[string]*utils.DistributedMutex),
	}
}

func (s *SmsSender)GetHasSendTimes() int64{
	return s.hasSendTimes
}


func (s *SmsSender) Flush() error {
	conn, err := redis.Dial("tcp", s.redisHost)
	if err != nil {
		logs.Error(err)
		return err
	}
	if _, err := conn.Do("flushall"); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s *SmsSender) Send(userId string) *model.ResultOfSend {

	// 双重判断，初始化用户的分布式锁
	if s.userDistributedMutex[userId] == nil {
		s.senderDistributeMutex.Lock()
		if s.userDistributedMutex[userId] == nil {
			s.userDistributedMutex[userId] = utils.NewDistributedMutex(
				s.redisHost,
				fmt.Sprintf("userLock:%s", userId),
				1,
			)
		}
		s.senderDistributeMutex.Unlock()
	}

	s.userDistributedMutex[userId].Lock()
	defer s.userDistributedMutex[userId].Unlock()

	userCertification, err := s.userCertificationStorage.GetUserCertification(userId)
	if err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       "消息发送失败，凭证获取失败",
		}
	}
	if userCertification == nil {
		userCertification = &model.Certification{
			RemainSendTimes: s.sendTimesPerTerm,
		}
	}
	if time.Now().Sub(userCertification.LastSendTime) >= s.term {
		userCertification.RemainSendTimes = s.sendTimesPerTerm
	}

	if userCertification.RemainSendTimes == 0 {
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       "发送过快，请稍后发送",
		}
	}

	// s.hasSendTimes++	// 这样会出现并发问题，破坏正确性，因为可能多个用户位于这(单个用户由于加锁了，所以不会同时位于这)
	// 这里发送短信的操作，可以交予消息队列
	atomic.AddInt64(&s.hasSendTimes, 1)
	// logs.Info("SMS 已发送: %d 次", atomic.AddInt64(&s.hasSendTimes, 1)) // 使用原子操作保证正确性


	userCertification.RemainSendTimes--
	userCertification.LastSendTime = time.Now()
	if err := s.userCertificationStorage.SetUserCertification(userId, userCertification); err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: true,
			Msg:       "消息发送成功，但凭证保存失败",
		}
	}
	return &model.ResultOfSend{
		IsSuccess: true,
		Msg: fmt.Sprintf(
			"发送成功, 剩余发送次数: %d, 发送时间: %v,",
			userCertification.RemainSendTimes,
			userCertification.LastSendTime,
		),
	}
}
