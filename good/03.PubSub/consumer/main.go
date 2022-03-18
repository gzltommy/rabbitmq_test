package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 1、连接rabbitmq
	conn, err := amqp.Dial("amqp://zorro:zorro@192.168.24.147:5672/zorro_test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 通过协程创建5个消费者
	for i := 0; i < 5; i++ {
		go func(number int) {
			// 2、创建信道，通常一个消费者一个
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel")
			defer ch.Close()

			// 3、声明交换机
			err = ch.ExchangeDeclare(
				"测试广播交换机", // 交换机名，需要跟消息发送方保持一致
				"fanout",  // 交换机类型
				true,      // 是否持久化
				false,     // auto-deleted
				false,     // internal
				false,     // no-wait
				nil,       // arguments
			)
			failOnError(err, "Failed to declare an exchange")

			// 4、声明需要操作的队列
			q, err := ch.QueueDeclare(
				"",    // 队列名字，不填则随机生成一个（这种随机生成的队列名的队列，在使用结束后队列被自动删除）
				false, // 是否持久化队列
				false, // delete when unused
				true,  // exclusive
				false, // no-wait
				nil,   // arguments
			)
			failOnError(err, "Failed to declare a queue")

			//5、队列绑定指定的交换机
			err = ch.QueueBind(
				q.Name,    // 队列名
				"",        // 路由参数，fanout 类型交换机，自动忽略路由参数
				"测试广播交换机", // 交换机名字，需要跟消息发送端定义的交换器保持一致
				false,
				nil)
			failOnError(err, "Failed to bind a queue")

			// 6、创建消费者
			msgs, err := ch.Consume(
				q.Name, // 引用前面的队列名
				"",     // 消费者名字，不填自动生成一个
				true,   // 自动向队列确认消息已经处理
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			failOnError(err, "Failed to register a consumer")

			// 循环消费队列中的消息
			for d := range msgs {
				log.Printf("[消费者编号=%d] 接收消息:%s", number, d.Body)
			}
		}(i)
	}

	// 挂起主协程，避免程序退出
	forever := make(chan bool)
	<-forever
}
