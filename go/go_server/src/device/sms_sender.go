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
	// 用户的分布式锁的创建信息
	redisHostOfUserDistributedMutex string
	keyPrefixOfUserDistributedMutex string
	retryTimeOfUserDistributedMutex time.Duration

	// 用户与分布式锁的映射关系
	userDistributedMutex        map[string]*middleware.DistributedMutex

	// 用户与凭证的映射关系
	userCertificationStorage    *middleware.TemporaryDataStorage


	// sms 订单号生成器
	orderNumDistributedGenerator *middleware.OrderNumDistributedGenerator

	// sms 的分布式锁
	smsSenderDistributedMutex   *middleware.DistributedMutex

	// 每个用户最多在 term 期间发送多少次短信
	sendTimesPerTerm            int

	// term 的时长
	term                        time.Duration

	// sms 已发送的短信次数
	hasSendTimes                int64

	// sms 的消息队列
	mqSender                    *middleware.MqSender

	// sms 发送消息到消息队列的 topic
	mqTopic                     string

	// sms 的 Id
	smsSenderId int
}

func NewSmsSender(
	redisHostOfUserDistributedMutex string,
	keyPrefixOfUserDistributedMutex string,
	retryTimeOfUserDistributedMutex time.Duration,

	sendTimesPerTerm int,
	term time.Duration,
	userCertificationStorage *middleware.TemporaryDataStorage,
	mqSender *middleware.MqSender,
	mqTopic string,
	smsSenderId int,
	orderNumDistributedGenerator *middleware.OrderNumDistributedGenerator,
	smsSenderDistributedMutex *middleware.DistributedMutex,
) *SmsSender {
	return &SmsSender{
		keyPrefixOfUserDistributedMutex:keyPrefixOfUserDistributedMutex,
		redisHostOfUserDistributedMutex: redisHostOfUserDistributedMutex,
		smsSenderDistributedMutex: smsSenderDistributedMutex,
		sendTimesPerTerm:         sendTimesPerTerm,
		term:                     term,
		userCertificationStorage: userCertificationStorage,
		userDistributedMutex:     make(map[string]*middleware.DistributedMutex),
		mqSender:                 mqSender,
		mqTopic:                  mqTopic,
		smsSenderId:smsSenderId,
		orderNumDistributedGenerator:orderNumDistributedGenerator,
		retryTimeOfUserDistributedMutex:retryTimeOfUserDistributedMutex,
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
	conn, err := redis.Dial("tcp", s.redisHostOfUserDistributedMutex)
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
				s.redisHostOfUserDistributedMutex,
				fmt.Sprintf("%s:%s", s.keyPrefixOfUserDistributedMutex,message.SenderId),
				s.retryTimeOfUserDistributedMutex,
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
		Key:   sarama.StringEncoder(model.MqMessage),
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
