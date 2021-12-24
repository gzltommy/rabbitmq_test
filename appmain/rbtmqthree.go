package main

import (
	"rbtmq/rbtmqcs"
)

func main() {
	rabbitmq := rbtmqcs.NewRabbitMQSimple("queueone")
	rabbitmq.ConsumeSimple()
}
