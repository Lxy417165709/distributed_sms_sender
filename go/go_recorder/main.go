package main

import (
	"distributed/go_recorder/src/object"
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func main() {
	kafkaHosts := []string{
		"qfstudio.net:15000",
		//"120.26.162.39:15000",
		//"120.26.162.39:15001",
		//"120.26.162.39:15002",
		//"120.26.162.39:9092",
	}
	topic := "sms_sender"
	mqReceiver := object.NewMqReceiver(kafkaHosts)
	if err := mqReceiver.Receive(topic); err != nil {
		logs.Error(err)
		return
	}
}
