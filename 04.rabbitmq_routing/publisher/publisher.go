package main

import (
	"fmt"
	"rbtmq/rabbitmq"
	"strconv"
	"time"
)

func main() {
	//路由模式下通过 key 将队列绑定到交换机上 这个队列式内部自动生成的 不需要指定名称 用的时候传递 key 即可找到绑定的 queue 队列
	//比如下边 传递了两个 key 那么就会内部绑定到交换机上两个队列  消费者只需要传递 key 过去即可找到交换机上绑定好的队列里面的消息进行消费
	rabbitmqOne := rabbitmq.NewRabbitMQRouting("exchange_routing", "route_key_one")
	rabbitmqTwo := rabbitmq.NewRabbitMQRouting("exchange_routing", "route_key_two")
	for i := 0; i <= 100; i++ {
		rabbitmqOne.PublishRouting("hello route_key_one " + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("hello route_key_two " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
