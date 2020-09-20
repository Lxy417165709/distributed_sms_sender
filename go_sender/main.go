package main

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var senders = []string{
	"1","2","3","4","5","6","7","8","9","10",
	"11","21","31","41","51","61","71","81","91","101",
	"12","22","32","42","52","62","72","82","92","102",
	"10017","10015","10013","10011","10009","10007","10005","10003","10001",
	"10000","10016","10014","10012","10010","10008","10006","10004","10002",
}

var receivers = []string{
	"1","2","3","4","5","6","7","8","9","10",
	"11","21","31","41","51","61","71","81","91","101",
	"12","22","32","42","52","62","72","82","92","102",
}

var contents = []string{
	"a","b","c","d",
	"我爱吃饭","我想你了","你是猪",
	"gopher","testing",
	"鱼我所欲","giao giao giao!!!",
	"abcder",":|",":)",":(","我是李白","嘻嘻嘻",
}


type SmsRequester struct{
	client *http.Client
	smsSenderUrl string
	identityNum int
}

func NewSmsRequester(smsSenderUrl string,identityNum int) *SmsRequester{
	return &SmsRequester{
		client: &http.Client{Timeout: 20 * time.Second},
		smsSenderUrl: smsSenderUrl,
		identityNum: identityNum,
	}
}

func (s *SmsRequester) SendRepeated() {
	go func() {
		for {
			go func() {
				sender := senders[rand.Intn(len(senders))]
				receiver := receivers[rand.Intn(len(receivers))]
				result,err := s.SendRequest(&MessageForm{
					SenderId: sender ,
					ReceiverId: receiver,
					Content: contents[rand.Intn(len(contents))],
				})
				if err != nil{
					logs.Error(err)
					return
				}
				logs.Info("[%d] Requester, %s -> %s,返回结果: %s",
					s.identityNum,
					sender,
					receiver,
					string(result),
				)
			}()
			time.Sleep(time.Duration(rand.Intn(2000))*time.Millisecond)
		}
	}()
}

func (s *SmsRequester) SendRequest(form *MessageForm) ([]byte,error){
	result,err := post(s.client,s.smsSenderUrl,"application/json",form)
	return result,err
}

func post(client *http.Client,url string, contentType string, data interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logs.Error(err)
		}
	}()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return result, nil
}


type MessageForm struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	countOfSmsRequester := 18
	smsSenderUrl := "http://127.0.0.1:8080/send_msg"
	smsRequesters := make([]*SmsRequester,0)
	for i:=0;i<countOfSmsRequester;i++{
		smsRequesters = append(smsRequesters,NewSmsRequester(smsSenderUrl,i))
	}
	for _,smsRequester := range smsRequesters{
		smsRequester.SendRepeated()
	}
	select{}
}
