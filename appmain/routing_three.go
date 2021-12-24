//路由模式消费者二
package main

import "rbtmq/rbtmqcs"

func main() {
	rabbitmqTwo := rbtmqcs.NewRabbitMQRouting("exHxb", "xiaobai_two")
	rabbitmqTwo.RecieveRouting()
}
