package main

import (
	"flag"
	"github.com/go-ozzo/ozzo-config"
	"os"
	"os/signal"
	"runtime"
	"superconsumer/constant"
	"superconsumer/log"
	"superconsumer/queue"
	"superconsumer/rpc"
	"sync"
	"sync/atomic"
)

type stats struct {
	taskNum    uint32 `已处理的任务数量`
	failedNum  uint32 `处理失败的任务数量`
	successNum uint32 `处理成功的任务数量`
}

type App struct {
	C        *config.Config
	Sig      chan os.Signal
	mChannel chan queue.Message
	stats    stats
}

type Task struct {
	ServiceName string
	MethodName  string
	RpcGroup    string
}

var (
	app = App{
		C:     config.New(),
		Sig:   make(chan os.Signal, 1),
		stats: stats{},
	}
	rpcClient = make(map[string]*rpc.RpcClient)
	consumer  queue.QueueInterface
	taskList  = make(map[string][]Task)
	Wg        sync.WaitGroup
)

func init() {
	c := flag.String("c", "/etc/superconsumerr.json", "配置文件")
	flag.Parse()
	if err := app.C.Load(*c); err != nil {
		panic(err)
	}
	signal.Notify(app.Sig, os.Interrupt)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app.mChannel = make(chan queue.Message, app.C.GetInt("maxConcurrent")) //次数可以通过获取配置文件中的最大并发数
	//配置logger
	log.NewLogger(app.C)
	os.Stdout.WriteString("日志组件配置完成...\n")

	//初始化rpc客户端
	initRpcClient()
	os.Stdout.WriteString("任务处理客户端初始化完成...\n")

	//初始化任务
	initTask()
	os.Stdout.WriteString("任务列表初始化完成...\n")

	//初始化队列消费者
	initConsumer()
	os.Stdout.WriteString("队列驱动初始化完成...\n")
	Wg.Add(2)
	//启动队列监听
	os.Stdout.WriteString("正在监听队列...\n")
	go listen()
	os.Stdout.WriteString("任务处理准备就绪...\n")
	go processTask()
	Wg.Wait()
}

// @description 初始化rpc客户端
func initRpcClient() {
	var rpc_config_group map[string]rpc.RpcConfig
	app.C.Configure(&rpc_config_group, "rpc")
	if len(rpc_config_group) == 0 {
		panic("rpc配置信息不合法")
	}

	for key, val := range rpc_config_group {
		rpcClient[key] = rpc.NewRpcClient(val)
	}
}

// @description 初始化队列消费者
func initConsumer() {
	consumer = queue.NewQueue(app.C, app.mChannel)
}

// @description 创建任务
func initTask() {
	if err := app.C.Configure(&taskList, "topicGroup"); err != nil {
		panic(err)
	}

	//遍历所有task，校验RpcGroup是否合法
	for topic, tasks := range taskList {
		for key, task := range tasks {
			if res, ok := rpcClient[task.RpcGroup]; !ok || res == nil {
				panic(task)
			}

			//校验method service
			if task.ServiceName == "" {
				panic("task's serviceName is empty!")
			}

			if task.MethodName == "" {
				taskList[topic][key].MethodName = "handleQueueContent"
			}
		}
	}
}

//@description 启动队列监听
func listen() {
	defer Wg.Done()
	consumer.Listen()
}

// @description 任务处理
func processTask() {
	defer Wg.Done()
	var wg sync.WaitGroup
loop:
	for {
		select {
		case message, ok := <-app.mChannel:
			if !ok {
				break loop
			}
			app.stats.taskNum++ //任务出来数量累加器
			for _, task := range taskList[message.TopicName] {
				//发送http请求
				client := rpcClient[task.RpcGroup]
				wg.Add(1)
				go func() {
					wg.Done()
					ret := client.Do(task.ServiceName, task.MethodName, message)
					if ret == constant.SUCCESS {
						atomic.AddUint32(&app.stats.successNum, 1)
					} else {
						atomic.AddUint32(&app.stats.failedNum, 1)
					}
				}()
			}
		case <-app.Sig:
			log.Info(
				"application",
				"任务处理数量 %d 失败数量 %d 成功数量 %d",
				app.stats.taskNum, app.stats.failedNum, app.stats.successNum)
			break loop
		}
	}
	wg.Wait()
}
