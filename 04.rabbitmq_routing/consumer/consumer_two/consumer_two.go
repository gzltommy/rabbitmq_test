package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmqTwo := rabbitmq.NewRabbitMQRouting("exchange_routing", "route_key_two")
	rabbitmqTwo.ReceiveRouting()
}

//你还可以搞更多的消费者 代码都一样 消费者越多那么读取队列里面消息的速度也就越快
