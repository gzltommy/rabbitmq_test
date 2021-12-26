package main

import (
	"rbtmq/rabbitmq"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQRoutingTest("queue_one", "exchange_test", "route_key_common")
	rabbitmqOne.ReceiveRoutingTest()
}
