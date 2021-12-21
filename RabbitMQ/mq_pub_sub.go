package RabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//这里是订阅模式的相关代码。

//订阅模式需要用到 exchange。

// NewRabbitMqSubscription 获取订阅模式下的 rabbitmq 的实例
func NewRabbitMqSubscription(exchangeName string) *RabbitMQ {
	//创建 rabbitmq 实例
	rabbitmq := NewRabbitMQ("", exchangeName, "")

	//获取 connection
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "订阅模式连接rabbitmq失败。")

	//获取 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "订阅模式获取channel失败")
	return rabbitmq
}

// PublishSubscription 订阅模式发布消息
func (r *RabbitMQ) PublishSubscription(message string) {
	//第一步，尝试连接交换机
	err := r.channel.ExchangeDeclare(
		r.ExChange,
		"fanout", //这里一定要设计为"fanout"也就是广播类型。
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "订阅模式发布方法中尝试连接交换机失败。")

	// 第二步，发送消息
	err = r.channel.Publish(
		r.ExChange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// ConsumeSubscription 订阅模式消费者
func (r *RabbitMQ) ConsumeSubscription() {
	//第一步，试探性创建交换机 exchange
	err := r.channel.ExchangeDeclare(
		r.ExChange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "订阅模式消费方法中创建交换机失败。")

	//第二步，试探性创建队列 queue
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "订阅模式消费方法中创建创建队列失败。")

	//第三步，绑定队列到交换机中
	err = r.channel.QueueBind(
		q.Name,
		"", // 在 pub/sub 模式下 key 要为空
		r.ExChange,
		false,
		nil,
	)

	//第四步，消费消息
	msgChannel, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgChannel {
			fmt.Printf("小杜同学写的订阅模式收到的消息：%s\n", d.Body)
		}
	}()

	fmt.Println("订阅模式消费者已开启，退出请按 CTRL+C\n")
	<-forever
}