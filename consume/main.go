// publish/main.go
package main

import (
	"fmt"
	"go_test/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMq.NewRabbitMQSimple("test_queue_name")
	rabbitmq.ConsumeSimple()
	fmt.Println("接收成功！")
}