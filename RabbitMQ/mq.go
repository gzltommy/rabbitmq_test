package RabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// 这里主要是 RabbitMQ 的一些信息。包括其结构体和函数。

// MQURL 连接信息
const MQURL = "amqp://guest:guest@192.168.24.147:5672/"

// RabbitMQ 结构体
type RabbitMQ struct {
	//连接
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列
	QueueName string
	//交换机名称
	ExChange string
	//绑定的key名称
	Key string
	//连接的信息，上面已经定义好了
	MqUrl string
}

// NewRabbitMQ 创建结构体实例，参数队列名称、交换机名称和 bind 的 key（也就是几个大写的，除去定义好的常量信息）
func NewRabbitMQ(queueName string, exChange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, ExChange: exChange, Key: key, MqUrl: MQURL}
}

// Destroy 关闭 conn 和 chanel 的方法
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// failOnErr 错误的函数处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		fmt.Printf("err是:%s,小杜同学手写的信息是:%s", err, message)
	}
}
