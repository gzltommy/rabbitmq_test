package RabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// rabbitmq 的路由模式。
// 主要特点不仅一个消息可以被多个消费者消费还可以由生产端指定消费者。
// 这里相对比订阅模式就多了一个 routingKey 的设计，也是通过这个来指定消费者的。
// 创建 exchange 的 kind 需要是"direct",不然就不是 roting 模式了。

// NewRabbitMqRouting 创建 rabbitmq 实例，这里有了 routingKey 为参数了。
func NewRabbitMqRouting(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)

	//获取connection
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "创建 rabbit 的路由实例的时候连接出现问题")

	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "创建 rabbitmq 的路由实例时获取channel出错")
	return rabbitmq
}

// PublishRouting 路由模式，产生消息。
func (r *RabbitMQ) PublishRouting(message string) {
	//第一步，尝试创建交换机，与 pub/sub 模式不同的是这里的 kind 需要是 direct
	err := r.channel.ExchangeDeclare(r.ExChange, "direct", true, false, false, false, nil)
	r.failOnErr(err, "路由模式，尝试创建交换机失败")

	//第二步，发送消息
	err = r.channel.Publish(
		r.ExChange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// ConsumerRouting 路由模式，消费消息。
func (r *RabbitMQ) ConsumerRouting() {
	//第一步，尝试创建交换机，注意这里的交换机类型与发布订阅模式不同，这里的是 direct
	err := r.channel.ExchangeDeclare(
		r.ExChange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "路由模式，创建交换机失败。")

	//第二步，尝试创建队列,注意这里队列名称不用写，这样就会随机产生队列名称
	q, err := r.channel.QueueDeclare(
		"", //随机产生队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "路由模式，创建队列失败。")

	//第三步，绑定队列到 exchange 中
	err = r.channel.QueueBind(q.Name, r.Key, r.ExChange, false, nil)

	//第四步，消费消息。
	messages, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			fmt.Printf("小杜同学写的路由模式(routing模式)收到消息为：%s。\n", d.Body)
		}
	}()
	<-forever
}
