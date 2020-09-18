package utils

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestDistributedCache(t *testing.T) {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	redisHosts := []string{
		"120.26.162.39:20000",
		"120.26.162.39:20001",
		"120.26.162.39:20002",
		"120.26.162.39:20003",
	}
	hashLoop := NewHashLoop(30, map[int64]*redis.Pool{
		0:  GetRedisConnPool(redisHosts[0]),
		8:  GetRedisConnPool(redisHosts[1]),
		16: GetRedisConnPool(redisHosts[2]),
		24: GetRedisConnPool(redisHosts[3]),
	})
	distributedCache := NewDistributedCache(hashLoop)
	keyValues := map[string]string{
		"":           "a",
		"klaskdjkaa": "a",
		"wqdwq":      "afcsafa",
		"d":          "afasfa",
		"dwqq":       "wqe q qw qwa",
		"ewqewqe":    "aeq weq qqe",
		"ksaiiwiqeqw12474 wq01e wqlklpowqpo ":"eklwqlk kcxvkskkloqw owq oeoqwoeqwu",
		"ksaiiwiqeqw101e wqlklpowqpo ":"eklwqlk kcxvkskklw2wq25 5w  owq oeoqwoeqwu",
		"ksaiiwiqeqw12474 1e wqlklpowqpo ":"eklwqlk kcxvk`odoqp pwskkloqw owq oeoqwoeqwu",
		"ksaiiwiqe74 wq01ewqlklpowqpo ":"eklwqlk kcxvksgfq qwq122cxz$*@(9$(%()kkloqw owq oeoqwoeqwu",
	}
	for key, value := range keyValues {
		if err := distributedCache.Set(key, value); err != nil {
			logs.Error(err)
			return
		}
		reply, err := distributedCache.Get(key)
		if err != nil {
			logs.Error(err)
			return
		}
		replyString, err := redis.String(reply, err)
		if err != nil {
			logs.Error(err)
			return
		}
		if replyString != fmt.Sprintf(`"%s"`,value) {
			logs.Error("key(%s)'s value is expected for %s, not %s",
				key,
				fmt.Sprintf(`"%s"`,value),
				replyString,
			)
			return
		}
	}
}
