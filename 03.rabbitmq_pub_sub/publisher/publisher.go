package main

import (
	"fmt"
	"rbtmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("exchange_pub_sub")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生成第" + strconv.Itoa(i) + "条数据")
		fmt.Println("订阅模式生成第" + strconv.Itoa(i) + "条数据")
		time.Sleep(1 * time.Second)
	}
}
