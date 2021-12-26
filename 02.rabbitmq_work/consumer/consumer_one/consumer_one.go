package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple("queue_work")
	rabbitmq.ConsumeSimple()
}
