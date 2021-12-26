package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple("queue_simple")
	rabbitmq.ConsumeSimple()
}
