## 简介
这是一个分布式的短信发送器(短信发送没实现)。实现了一段时间内，同一用户最多发送 `max` 次短信的限制。

## 技术栈
- [ ] 前端
	- `Nginx`
- [ ] 后端
	- `Go`
- [ ] 存储
	- `Redis`
	- `Mysql`
- [ ] 中间件
	- `Docker`
	- `Kafka`
	- `Grafana`

## 特点
1. 同一用户可位于不同机器。
2. 实现 `Docker` 实现所有组件的容器化部署。
3. 后端以 `Go` 集群形态为用户服务。
4. 使用 `Redis` 实现分布式锁。
5. 使用 `Redis`集群  + 一致性哈希算法 实现分布式缓存。
6. 使用 `Nginx` 实现负载均衡。
7. 使用 `Kafka` 实现消息队列。
8. 使用 `Grafana` 实现 `Mysql` 数据的图表化。
9. 代码优美、测试较完备。  `:)`


## 相关信息
1. [grafana 数据可视化在线观看](http://120.26.162.39:35000/d/qhpaOPdMz/invoke_situation?orgId=1&from=1600560096000&to=1600563340000)
2. [架构](https://www.processon.com/view/link/5f64edc66376894e3278ceba)

