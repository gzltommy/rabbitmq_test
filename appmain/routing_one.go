//路由模式生产者
package main

import (
	"fmt"
	"rbtmq/rbtmqcs"
	"strconv"
	"time"
)

func main() {
	//路由模式下通过key将队列绑定到交换机上 这个队列是内部自动生成的 不需要指定名称 用的时候传递key即可找到绑定的queue队列
	//比如下边 传递了两个key那么就会内部绑定到交换机上两个队列  消费者只需要传递key过去即可找到交换机上绑定好的队列里面的消息进行消费
	rabbitmqOne := rbtmqcs.NewRabbitMQRouting("exHxb", "xiaobai_one")
	rabbitmqTwo := rbtmqcs.NewRabbitMQRouting("exHxb", "xiaobai_two")
	for i := 0; i <= 10; i++ {
		rabbitmqOne.PublishRouting("hello xiaobai one" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("hello xiaobai two" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
