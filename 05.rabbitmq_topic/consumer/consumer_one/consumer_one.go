package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	//# 号表示匹配多个单词 也就是读取 exchange_topic 交换机里面所有队列的消息
	rabbitmq := rabbitmq.NewRabbitMQTopic("exchange_topic", "#")
	rabbitmq.RecieveTopic()
}
