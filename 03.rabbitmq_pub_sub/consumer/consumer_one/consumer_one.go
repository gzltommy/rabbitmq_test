package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("exchange_pub_sub")
	rabbitmq.ReceiveSub()
}
