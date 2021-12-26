package main

import (
	"fmt"
	"rbtmq/rabbitmq"
)

func main() {
	rabitmq := rabbitmq.NewRabbitMQSimple("queue_simple")
	rabitmq.PublishSimple("hello tommy!！")
	fmt.Println("发送成功!")
}
