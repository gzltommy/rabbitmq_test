package main

import "rbtmq/rbtmqcs"

func main() {
	rabbitmq := rbtmqcs.NewRabbitMQSimple("queuetwo")
	rabbitmq.ConsumeSimple()
}
