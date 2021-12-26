package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	//这里只是匹配到了 tommy.后边只能是一个单词的 key 通过这个key去找绑定到交换机上的相应的队列
	rabbitmq := rabbitmq.NewRabbitMQTopic("exchange_topic", "tommy.*.cs")
	rabbitmq.RecieveTopic()
}
