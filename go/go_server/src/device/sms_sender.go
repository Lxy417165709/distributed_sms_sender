package device

import (
	"distributed/sms/src/middleware"
	"distributed/sms/src/model"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"sync/atomic"
	"time"
)

type SmsSender struct {
	redisHostOfDistributedMutex string
	orderNumDistributedGenerator *middleware.OrderNumDistributedGenerator
	userCertificationStorage    *middleware.TemporaryDataStorage
	smsSenderDistributedMutex   *middleware.DistributedMutex
	userDistributedMutex        map[string]*middleware.DistributedMutex
	sendTimesPerTerm            int
	term                        time.Duration
	hasSendTimes                int64
	mqSender                    *middleware.MqSender
	mqTopic                     string
	smsSenderId int
}

func NewSmsSender(
	redisHostOfDistributedMutex string,
	sendTimesPerTerm int,
	term time.Duration,
	userCertificationStorage *middleware.TemporaryDataStorage,
	mqSender *middleware.MqSender,
	mqTopic string,
	smsSenderId int,
	orderNumDistributedGenerator *middleware.OrderNumDistributedGenerator,
) *SmsSender {
	return &SmsSender{
		redisHostOfDistributedMutex: redisHostOfDistributedMutex,
		smsSenderDistributedMutex: middleware.NewDistributedMutex(
			redisHostOfDistributedMutex,
			"smsSenderLock",
			1),
		sendTimesPerTerm:         sendTimesPerTerm,
		term:                     term,
		userCertificationStorage: userCertificationStorage,
		userDistributedMutex:     make(map[string]*middleware.DistributedMutex),
		mqSender:                 mqSender,
		mqTopic:                  mqTopic,
		smsSenderId:smsSenderId,
		orderNumDistributedGenerator:orderNumDistributedGenerator,
	}
}

func (s *SmsSender) GetHasSendTimes() int64 {
	return s.hasSendTimes
}

func (s *SmsSender) Flush() error {
	if err := s.userCertificationStorage.Flush(); err != nil {
		logs.Error(err)
		return err
	}
	conn, err := redis.Dial("tcp", s.redisHostOfDistributedMutex)
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

func (s *SmsSender) Send(message *model.Message) *model.ResultOfSend {
	// 双重判断，初始化用户的分布式锁
	if s.userDistributedMutex[message.SenderId] == nil {
		s.smsSenderDistributedMutex.Lock()
		if s.userDistributedMutex[message.SenderId] == nil {
			s.userDistributedMutex[message.SenderId] = middleware.NewDistributedMutex(
				s.redisHostOfDistributedMutex,
				fmt.Sprintf("userLock:%s", message.SenderId),
				1,
			)
		}
		s.smsSenderDistributedMutex.Unlock()
	}

	// 这里进行加锁，防止发送消息的同一个用户同时位于临界区
	s.userDistributedMutex[message.SenderId].Lock()
	defer s.userDistributedMutex[message.SenderId].Unlock()

	userCertification, err := s.getUserNewestCertification(message.SenderId)
	if err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       "消息发送失败，凭证获取失败",
		}
	}

	if s.hasUserCertificationExpired(userCertification) {
		userCertification.RemainSendTimes = s.sendTimesPerTerm
	}

	if userCertification.RemainSendTimes == 0 {
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       "发送过快，请稍后发送",
		}
	}

	message.OrderNum,err = s.orderNumDistributedGenerator.GetOrderNum()
	if err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       err.Error(),
		}
	}
	message.SmsSenderId = s.smsSenderId


	messageJson, err := json.Marshal(message)
	if err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       err.Error(),
		}
	}

	// 将存储任务交予消息队列
	if err := s.mqSender.Send(&sarama.ProducerMessage{
		Topic: s.mqTopic,
		Key:   sarama.StringEncoder("message"),
		Value: sarama.StringEncoder(messageJson),
	}); err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: false,
			Msg:       err.Error(),
		}
	}

	// 执行发送短信
	atomic.AddInt64(&s.hasSendTimes, 1)
	userCertification.RemainSendTimes--
	userCertification.LastSendTime = time.Now()

	if err := s.userCertificationStorage.SetUserCertification(message.SenderId, userCertification); err != nil {
		logs.Error(err)
		return &model.ResultOfSend{
			IsSuccess: true,
			Msg:       "短信发送成功，但凭证保存失败",
		}
	}
	return &model.ResultOfSend{
		IsSuccess: true,
		Msg: fmt.Sprintf(
			"发送成功, 发送时间: %v, 剩余发送次数: %d",
			userCertification.LastSendTime,
			userCertification.RemainSendTimes,
		),
	}
}

func (s *SmsSender) hasUserCertificationExpired(userCertification *model.Certification) bool {
	return time.Now().Sub(userCertification.LastSendTime) >= s.term
}

func (s *SmsSender) getUserNewestCertification(userId string) (*model.Certification, error) {
	userCertification, err := s.userCertificationStorage.GetUserCertification(userId)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	if userCertification == nil {
		userCertification = &model.Certification{
			RemainSendTimes: s.sendTimesPerTerm,
		}
	}
	return userCertification, nil
}
