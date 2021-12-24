package main

import (
	"fmt"
	"rbtmq/rbtmqcs"
)

func main() {
	rabitmq := rbtmqcs.NewRabbitMQSimple("queueone")
	rabitmq.PublishSimple("hello huxiaobai12345!！！")
	fmt.Println("发送成功!")
}
