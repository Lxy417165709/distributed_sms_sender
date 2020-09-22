## 简介
这是一个分布式的短信发送器(短信发送没实现)。实现了一段时间内，同一用户最多发送 `max` 次短信的限制。

## 技术栈
- [ ] 后端
	- `Go`
- [ ] 存储
	- `Redis`
	- `Mysql`
- [ ] 中间件
	- `Docker`
	- `Kafka`
	- `Nginx`
	- `Grafana`
- [ ] 其他
	- `yaml`
	- `shell`
	- `lua`

## 特点
1. 同一用户可位于不同机器。
2. 实现 `Docker` 实现所有组件的容器化部署。
3. 后端以 `Go` 集群形态为用户服务。
4. 使用 `Redis` 实现分布式锁。
5. 使用 `Redis`集群  + 一致性哈希算法 实现分布式缓存。
6. 使用 `Nginx` 实现负载均衡。
7. 使用 `Kafka` 实现消息队列。
8. 使用 `Grafana` 实现 `Mysql` 数据的图表化。
9. 使用 `lua` + `redis` 实现分布式生成唯一订单号。
10. 代码优美、测试较完备。  `:)`


## 工具
1. 两台阿里云服务器(**感谢宏宏的阿里云服务器~**)

## 相关信息
1. [grafana 数据可视化在线观看](http://120.26.162.39:35000/d/qhpaOPdMz/message_dashboard?orgId=1&refresh=10s)
2. [架构](https://www.processon.com/view/link/5f64edc66376894e3278ceba)

## 学习资料
1. [sarama 文档](http://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/go%E6%93%8D%E4%BD%9Ckafka/%E6%93%8D%E4%BD%9Ckafka.html)


