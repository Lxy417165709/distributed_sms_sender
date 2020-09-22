package device

import (
	"distributed/sender/src/model"
	"distributed/sender/src/utils"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"net/http"
	"time"
)

type SmsRequester struct {
	client       *http.Client
	smsSenderUrl string
	identityNum  int
}

func NewSmsRequester(smsSenderUrl string, identityNum int) *SmsRequester {
	return &SmsRequester{
		client:       &http.Client{Timeout: 20 * time.Second},
		smsSenderUrl: smsSenderUrl,
		identityNum:  identityNum,
	}
}

func (s *SmsRequester) SendRepeatedly(senders, receivers, contents []string) {
	go func() {
		for {
			go func() {
				sender := senders[rand.Intn(len(senders))]
				receiver := receivers[rand.Intn(len(receivers))]
				result, err := s.SendRequest(&model.MessageForm{
					SenderId:   sender,
					ReceiverId: receiver,
					Content:    contents[rand.Intn(len(contents))],
				})
				if err != nil {
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
			time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
		}
	}()
}

func (s *SmsRequester) SendRequest(form *model.MessageForm) ([]byte, error) {
	result, err := utils.Post(s.client, s.smsSenderUrl, "application/json", form)
	return result, err
}
