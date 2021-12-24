//work模式下的生产者
package main

import (
	"fmt"
	"rbtmq/rbtmqcs"
	"strconv"
)

func main() {
	rabbitmq := rbtmqcs.NewRabbitMQSimple("queuetwo")
	for i := 0; i <= 1000; i++ {
		//strconv.Itoa(i) 科普一下 可以将整形转换成字符串型的数字
		rabbitmq.PublishSimple("hello xiaobai" + strconv.Itoa(i))
		//time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
