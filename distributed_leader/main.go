package main

// 导入所需的库
import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"strings"
)

// 可以设置一些参数，比如节点ID
var (
	nodeID    = flag.Int("id", 0, "node ID")
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	electName = flag.String("name", "my-test-elect", "election name")
)

func main() {
	flag.Parse()

	// 将 etcd 的地址解析成 slice of string
	endpoints := strings.Split(*addr, ",")

	// 生成一个 etcd 的 clien
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建 session，如果程序宕机导致 session 断掉，etcd 能检测到
	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// 生成一个选举对象。下面主要使用它进行选举和查询等操作
	// 另一个方法 ResumeElection 可以使用既有的 leader 初始化 Election
	e1 := concurrency.NewElection(session, *electName)

	// 从命令行读取命令
	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		action := consoleScanner.Text()
		switch action {
		case "elect": // 选举命令
			go elect(e1, *electName)
		case "proclaim": // 只更新 leader 的 value
			proclaim(e1, *electName)
		case "resign": // 辞去 leader，重新选举
			resign(e1, *electName)
		case "watch": // 监控 leader 的变动
			go watch(e1, *electName)
		case "query": // 查询当前的 leader
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}

var count int

// 选主
func elect(e1 *concurrency.Election, electName string) {
	log.Println("acampaigning for ID:", *nodeID)

	// 调用 Campaign 方法选主，主的值为 value-<主节点ID>-<count>
	if err := e1.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("campaigned for ID:", *nodeID)
	count++
}

// 为主设置新值
func proclaim(e1 *concurrency.Election, electName string) {
	log.Println("proclaiming for ID:", *nodeID)

	// 调用 Proclaim 方法设置新值，新值为 value-<主节点ID>-<count>
	if err := e1.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("proclaimed for ID:", *nodeID)
	count++
}

// 重新选主，有可能另外一个节点被选为了主
func resign(e1 *concurrency.Election, electName string) {
	log.Println("resigning for ID:", *nodeID)

	// 调用 Resign 重新选主
	if err := e1.Resign(context.TODO()); err != nil {
		log.Println(err)
	}
	log.Println("resigned for ID:", *nodeID)
}

// 查询主的信息
func query(e1 *concurrency.Election, electName string) {
	// 调用 Leader 返回主的信息，包括 key 和 value 等信息
	resp, err := e1.Leader(context.Background())
	if err != nil {
		log.Printf("failed to get the current leader: %v", err)
	}
	log.Println("current leader:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
}

// 可以直接查询主的 rev 信息
func rev(e1 *concurrency.Election, electName string) {
	rev := e1.Rev()
	log.Println("current rev:", rev)
}

func watch(e1 *concurrency.Election, electName string) {
	ch := e1.Observe(context.TODO())

	log.Println("start to watch for ID:", *nodeID)
	for i := 0; i < 10; i++ {
		resp := <-ch
		log.Println("leader changed to", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
	}
}
