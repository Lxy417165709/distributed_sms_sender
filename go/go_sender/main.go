package main

import (
	"distributed/sender/src/device"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"time"
)

const (
	countOfSmsRequester = 30
	smsSenderUrl        = "http://127.0.0.1:8080/send_msg"
)

var (
	senders = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "21", "31", "41", "51", "61", "71", "81", "91", "101",
		"12", "22", "32", "42", "52", "62", "72", "82", "92", "102",
		"10017", "10015", "10013", "10011", "10009", "10007", "10005", "10003", "10001",
		"10000", "10016", "10014", "10012", "10010", "10008", "10006", "10004", "10002",
	}
	receivers = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "21", "31", "41", "51", "61", "71", "81", "91", "101",
		"12", "22", "32", "42", "52", "62", "72", "82", "92", "102",
	}
	contents = []string{
		"a", "b", "c", "d",
		"我爱吃饭", "我想你了", "你是猪",
		"gopher", "testing",
		"鱼我所欲", "giao giao giao!!!",
		"abcder", ":|", ":)", ":(", "我是李白", "嘻嘻嘻",
	}
)

func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	smsRequesters := make([]*device.SmsRequester, 0)
	for i := 0; i < countOfSmsRequester; i++ {
		smsRequesters = append(smsRequesters, device.NewSmsRequester(smsSenderUrl, i))
	}
	for _, smsRequester := range smsRequesters {
		smsRequester.SendRepeatedly(senders, receivers, contents)
	}
	select {}
}
