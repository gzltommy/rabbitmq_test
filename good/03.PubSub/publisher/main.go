package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 1、连接 rabbitmq
	conn, err := amqp.Dial("amqp://zorro:zorro@192.168.24.147:5672/zorro_test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	///2、创建信道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//3、声明交换机
	err = ch.ExchangeDeclare(
		"测试广播交换机", // 交换机名字
		"fanout",  // 交换机类型，fanout 发布订阅模式
		true,      // 是否持久化
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 4、推送消息
	for i := 0; i < 6; i++ {
		body := fmt.Sprintf("Hello! -- %d", i) // 消息内容
		err = ch.Publish(
			"测试广播交换机", // exchange（交换机名字，跟前面声明对应）
			"",        // 路由参数，fanout 类型交换机，自动忽略路由参数，填了也没用。
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain", // 消息内容类型，这里是普通文本
				Body:        []byte(body), // 消息内容
			})
		log.Printf("发送内容 %s", body)
	}

}
