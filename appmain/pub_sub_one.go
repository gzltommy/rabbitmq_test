//订阅模式下的生产者
package main

import (
	"fmt"
	"rbtmq/rbtmqcs"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rbtmqcs.NewRabbitMQPubSub("newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生成第" + strconv.Itoa(i) + "条" + "数据")
		fmt.Println("订阅模式生成第" + strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}
}
