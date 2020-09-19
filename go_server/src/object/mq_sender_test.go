package object

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

func init(){
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func TestMqSender(t *testing.T) {
	kafkaHosts := []string{"120.26.162.39:9092"}
	topic := "test_topic"
	for i:=0;;i++{
		mqSender := NewMqSender(kafkaHosts)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("key_%d",i)),
			Value: sarama.StringEncoder(fmt.Sprintf("value_%d",i)),
		}
		logs.Info("Sending { key: %s, value: %s } success", msg.Key, msg.Value)
		if err := mqSender.Send(msg); err != nil {
			logs.Error(err)
			return
		}
		time.Sleep(1*time.Second)
	}
}
