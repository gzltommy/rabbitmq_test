package main

import (
	"fmt"
	"rbtmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQTopic("exchange_topic", "tommy.one")
	rabbitmqTwo := rabbitmq.NewRabbitMQTopic("exchange_topic", "tommy.two.cs")
	for i := 0; i <= 10; i++ {
		rabbitmqOne.PublishTopic("hello tommy one" + strconv.Itoa(i))
		rabbitmqTwo.PublishTopic("hello tommy two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
