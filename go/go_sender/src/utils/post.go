package utils

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

func Post(client *http.Client, url string, contentType string, data interface{}) ([]byte, error) {
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
