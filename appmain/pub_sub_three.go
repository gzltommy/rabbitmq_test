//订阅模式下的消费者二
package main

import "rbtmq/rbtmqcs"

func main() {
	rabbitmq := rbtmqcs.NewRabbitMQPubSub("newProduct")
	rabbitmq.RecieveSub()
}
