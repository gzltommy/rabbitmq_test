package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQRouting("exchange_routing", "route_key_one")
	rabbitmqOne.ReceiveRouting()
}
