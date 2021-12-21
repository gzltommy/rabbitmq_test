// publish/main.go
package main

import (
	"fmt"
	"go_test/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMq.NewRabbitMQSimple("test_queue_name")
	rabbitmq.PublishSimple("他是客，你是心上人。 ---来自simple模式")
	fmt.Println("发送成功！")
}