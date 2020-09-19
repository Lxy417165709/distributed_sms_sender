package object

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type MqReceiver struct {
	kafkaConsumer sarama.Consumer
}

func NewMqReceiver(kafkaHosts []string) *MqReceiver {
	consumer, err := sarama.NewConsumer(
		kafkaHosts,
		nil,
	)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	return &MqReceiver{
		kafkaConsumer: consumer,
	}
}

func (m *MqReceiver) Receive(topic string) error{
	partitions, err := m.kafkaConsumer.Partitions(topic)
	if err != nil {
		logs.Error(err)
		return err
	}
	for _, partition := range partitions {
		partitionConsumer, err := m.kafkaConsumer.ConsumePartition(
			topic,
			partition,
			sarama.OffsetOldest,
		)
		if err != nil {
			logs.Error(err)
			return err
		}
		go func(partitionConsumer sarama.PartitionConsumer) {
			for {
				select {
				case m := <-partitionConsumer.Messages():
					logs.Info(
						"Receiving { key: %s, value: %s, offset: %d } success",
						m.Key,
						m.Value,
						m.Offset,
					)
				}
			}
		}(partitionConsumer)
	}
	select {}
}
