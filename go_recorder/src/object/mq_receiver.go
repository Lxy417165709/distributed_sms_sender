package object

import (
	"distributed/go_recorder/src/model"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)


type MqReceiver struct {
	kafkaConsumer sarama.Consumer
	db *gorm.DB
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
		db: getDB(),
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
				case msg := <-partitionConsumer.Messages():
					logs.Info(
						"Receiving { key: %s, value: %s, offset: %d } success",
						msg.Key,
						msg.Value,
						msg.Offset,
					)
					if !m.db.HasTable(&model.InvokeSituation{}) {
						m.db.CreateTable(&model.InvokeSituation{})
					}
					m.db.Create(&model.InvokeSituation{
						UserId: string(msg.Key),
						InvokeTime:time.Now(),
					})
					logs.Info(
						"saving { key: %s, value: %s, offset: %d } success",
						msg.Key,
						msg.Value,
						msg.Offset,
					)
				}
			}
		}(partitionConsumer)
	}
	select {}
}





func getDB() *gorm.DB {
	user := "root"
	password := "123456"
	host:="120.26.162.39:40000"
	dbName := "mysql"
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user,
			password,
			host,
			dbName,
		),
	)
	if err != nil {
		panic(err)
	}
	return db
}
