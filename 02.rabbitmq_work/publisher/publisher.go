package main

import (
	"fmt"
	"rbtmq/rabbitmq"
	"strconv"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple("queue_work")
	for i := 0; i <= 1000; i++ {
		//strconv.Itoa(i) 科普一下 可以将整形转换成字符串型的数字
		rabbitmq.PublishSimple("tommy 请说：" + strconv.Itoa(i))
		//time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
