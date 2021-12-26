package main

import (
	"fmt"
	"rbtmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmqV := rabbitmq.NewRabbitMQRoutingTest("", "exchange_test", "route_key_common")
	for i := 0; i <= 100; i++ {
		rabbitmqV.PublishRoutingTest("one key multiple queue test " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
