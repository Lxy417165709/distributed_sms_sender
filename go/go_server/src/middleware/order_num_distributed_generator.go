package middleware

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

type OrderNumDistributedGenerator struct {
	connPool    *redis.Pool
	orderNumKey string
}

func NewOrderNumDistributedGenerator(connPool *redis.Pool, orderNumKey string) *OrderNumDistributedGenerator {
	return &OrderNumDistributedGenerator{
		connPool:    connPool,
		orderNumKey: orderNumKey,
	}
}

func (o *OrderNumDistributedGenerator) GetOrderNum() (int, error) {
	const incrScript = `
		local num = redis.call('GET', KEYS[1]);  
		
		if not num then
			redis.call('SET',KEYS[1], 0); 
			return 0;	
		else
			local res = num + 1; 
			redis.call('SET',KEYS[1], res); 
			return res;
		end
	`

	lua := redis.NewScript(1, incrScript)
	redisConn := o.connPool.Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			logs.Error(err)
		}
	}()
	orderNum, err := redis.Int(lua.Do(redisConn, o.orderNumKey))
	if err != nil {
		logs.Error(err)
		return 0, err
	}
	return orderNum, nil
}
