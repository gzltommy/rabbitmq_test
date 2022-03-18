package main

// 导入包
import (
	"github.com/streadway/amqp"
	"log"
)

// 错误处理
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 1、连接 RabbitMQ
	conn, err := amqp.Dial("amqp://zorro:zorro@192.168.24.147:5672/zorro_test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2、创建信道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//3、声明要操作的队列
	q, err := ch.QueueDeclare(
		"hello", // 队列名需要跟发送消息的队列名保持一致
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 创建消息消费者
	msgChan, err := ch.Consume(
		q.Name, // 队列名
		"",     // 消费者名字，不填，则自动生成一个唯一ID
		true,   // 是否自动提交消息，即自动告诉rabbitmq消息已经处理成功。
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// 循环拉取队列中的消息
	for d := range msgChan {
		// 打印消息内容
		log.Printf("Received a message: %s", d.Body)
	}
}
