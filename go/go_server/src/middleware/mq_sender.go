package middleware

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type MqSender struct {
	kafkaProducer sarama.SyncProducer
}

func NewMqSender(kafkaHosts []string) *MqSender {
	kafkaConf := sarama.NewConfig()
	kafkaConf.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConf.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(
		kafkaHosts,
		kafkaConf,
	)
	if err!=nil{
		logs.Error(err)
		panic(err)
	}
	return &MqSender{
		kafkaProducer: producer,
	}
}


func (m *MqSender) Send(msg *sarama.ProducerMessage) error {
	_, _, err := m.kafkaProducer.SendMessage(msg)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}
