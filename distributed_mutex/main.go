package main

import (
	"context"
	"flag"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	lockName = flag.String("name", "my-test-lock", "lock name")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	// etcd 地址
	endpoints := strings.Split(*addr, ",")

	// 生成一个 etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	useLock(cli) // 测试锁
}

func useLock(cli *clientv3.Client) {
	// 为锁生成 session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	// 得到一个分布式锁
	locker := concurrency.NewLocker(s1, *lockName)

	// 请求锁
	log.Println("acquiring lock")
	locker.Lock()
	log.Println("acquired lock")

	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	// 释放锁
	locker.Unlock()
	log.Println("released lock")
}

func useMutex(cli *clientv3.Client) {
	// 为锁生成 session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	m1 := concurrency.NewMutex(s1, *lockName)

	//在请求锁之前查询 key
	log.Printf("before acquiring. key: %s", m1.Key())

	// 请求锁
	log.Println("acquiring lock")
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Printf("acquired lock. key: %s", m1.Key())

	//等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	// 释放锁
	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("released lock")
}
