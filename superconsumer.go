package main

import (
	"flag"
	"github.com/go-ozzo/ozzo-config"
	"superconsume/rpc"
	"superconsume/log"
	"superconsume/queue"
	"os"
	"sync"
	"sync/atomic"
	"superconsume/constant"
	"runtime"
	"os/signal"
)

type stats struct {
	taskNum uint32 `已处理的任务数量`
	failedNum uint32 `处理失败的数量`
	successNum uint32 `处理成功的数量`
}

type App struct {
	C *config.Config
	ConfigPath *string
	Sig chan os.Signal
	Wg *sync.WaitGroup
	mChannel chan queue.Message
	stats stats
}

type Task struct {
	ServiceName string
	MethodName string
	RpcGroup string
}

var (
	app = App{
		C: config.New(),
		ConfigPath:flag.String("c", "/etc/superconsumer.json", "配置文件"),
		Sig: make(chan os.Signal, 1),
		stats: stats{},
		Wg:&sync.WaitGroup{},
	}

	rpcClient = make(map[string]*rpc.RpcClient)

	consumer queue.QueueInterface

	taskList = make(map[string][]Task)

    Wg sync.WaitGroup
)

func main() {
	flag.Parse()

	//加载配置文件
	if err := app.C.Load(*app.ConfigPath); err != nil {
		panic(err)
	}
	maxConcurrent := app.C.GetInt("maxConcurrent")
	app.mChannel = make(chan queue.Message, maxConcurrent) //次数可以通过获取配置文件中的最大并发数

	signal.Notify(app.Sig,os.Interrupt)//监听退出信号、user1、user2
	initApp()
	runtime.GOMAXPROCS(runtime.NumCPU())
	Wg.Add(2)
	//启动队列监听
	go consumer.Listen()

	go processTask()
	Wg.Wait()
}

func initApp() {
	//配置logger
	log.NewLogger(app.C)

	//初始化rpc客户端
	initRpcClient()

	//初始化任务
	initTask()

	//初始化队列消费者
	initConsumer()
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

func initTask() {
	if err := app.C.Configure(&taskList, "topicGroup"); err != nil {
		panic(err)
	}

	//遍历所有task，校验RpcGroup是否合法
	for topic,tasks := range taskList {
		for key,task := range tasks {
			if res,ok := rpcClient[task.RpcGroup]; !ok || res == nil {
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

func processTask(){
	defer Wg.Done()

	var wg sync.WaitGroup
	loop:for {
		select {
		case message := <- app.mChannel:
			app.stats.taskNum ++ //任务出来数量累加器
			panic(message)
			for _,task := range taskList[message.TopicName] {
				//发送http请求
				client := rpcClient[task.RpcGroup]
				wg.Add(1)
				go func() {
					wg.Done()
					ret := client.Do(task.ServiceName, task.MethodName,message)
					if ret == constant.SUCCESS {
						atomic.AddUint32(&app.stats.successNum, 1)
					} else {
						atomic.AddUint32(&app.stats.failedNum, 1)
					}
				}()
			}
		case <- app.Sig:
			break loop
		default:
			break loop
		}
	}
}