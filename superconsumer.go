package main

import (
	"flag"
	"github.com/go-ozzo/ozzo-config"
	"os"
	"superconsumer/constant"
	"superconsumer/log"
	"superconsumer/queue"
	"superconsumer/rpc"
	"sync"
	"sync/atomic"
)

const (
	DefaultConfig = "./config/superconsumer.json"
)

type Task struct {
	ServiceName string
	MethodName  string
	RpcGroup    string
}

type stats struct {
	taskNum    uint32 `已处理的任务数量`
	failedNum  uint32 `处理失败的任务数量`
	successNum uint32 `处理成功的任务数量`
}

type App struct {
	sync.WaitGroup

	C        *config.Config
	mChannel chan queue.Message
	stats    stats

	rpcClient map[string]*rpc.RpcClient
	consumer  queue.QueueInterface
	taskList  map[string][]Task
}

// @description 初始化rpc客户端
func (app *App) initRpcClient() {
	var rpc_config_group map[string]rpc.RpcConfig
	app.C.Configure(&rpc_config_group, "rpc")
	if len(rpc_config_group) == 0 {
		panic("rpc配置信息不合法")
	}

	for key, val := range rpc_config_group {
		app.rpcClient[key] = rpc.NewRpcClient(val)
	}
}

// @description 初始化队列消费者
func (app *App) initConsumer() {
	app.consumer = queue.NewQueue(app.C, app.mChannel)
}

// @description 创建任务
func (app *App) initTask() {
	if err := app.C.Configure(&app.taskList, "topicGroup"); err != nil {
		panic(err)
	}

	//遍历所有task，校验RpcGroup是否合法
	for topic, tasks := range app.taskList {
		for key, task := range tasks {
			if res, ok := app.rpcClient[task.RpcGroup]; !ok || res == nil {
				panic(task)
			}

			//校验method service
			if task.ServiceName == "" {
				panic("task's serviceName is empty!")
			}
			if task.MethodName == "" {
				app.taskList[topic][key].MethodName = "handleQueueContent"
			}
		}
	}
}

//@description 启动队列监听
func (app *App) listen() {
	defer app.Done()
	app.consumer.Listen()
}

// @description 任务处理
func (app *App) processTask() {
	defer app.Done()

	var wg sync.WaitGroup
loop:
	for {
		select {
		case message, ok := <-app.mChannel:
			//写端已关闭。读端处理完channel中剩余的消息后在退出!
			if !ok {
				break loop
			}
			app.stats.taskNum++ //任务出来数量累加器
			for _, task := range app.taskList[message.TopicName] {
				wg.Add(1)

				//发送http请求
				client := app.rpcClient[task.RpcGroup]
				go func() {
					defer wg.Done()

					ret := client.Do(task.ServiceName, task.MethodName, message)
					if ret == constant.SUCCESS {
						atomic.AddUint32(&app.stats.successNum, 1)
					} else {
						atomic.AddUint32(&app.stats.failedNum, 1)
					}
				}()
			}
		}
	}
	wg.Wait()
}

func (app *App) Pretreatment() {

	app.LoadConfig()
	app.mChannel = make(chan queue.Message, app.C.GetInt("maxConcurrent")) //次数可以通过获取配置文件中的最大并发数

	//配置logger
	log.NewLogger(app.C)
	os.Stdout.WriteString("日志组件配置完成...\n")

	//初始化rpc客户端
	app.initRpcClient()
	os.Stdout.WriteString("任务处理客户端初始化完成...\n")

	//初始化任务
	app.initTask()
	os.Stdout.WriteString("任务列表初始化完成...\n")

	//初始化队列消费者
	app.initConsumer()
	os.Stdout.WriteString("队列驱动初始化完成...\n")
}

func (app *App) Run() {
	app.Add(2)

	//启动队列监听
	os.Stdout.WriteString("正在监听队列...\n")
	go app.listen()

	os.Stdout.WriteString("任务处理准备就绪...\n")
	go app.processTask()
}

func (app *App) LoadConfig() {
	c := flag.String("c", DefaultConfig, "配置文件")
	flag.Parse()
	if err := app.C.Load(*c); err != nil {
		panic(err)
	}
}

func main() {

	app := &App{
		C:         config.New(),
		stats:     stats{},
		rpcClient: make(map[string]*rpc.RpcClient),
		taskList:  make(map[string][]Task),
	}

	// Pretreatment用于预处理阶段：包括包含一些环境的初始化等等
	app.Pretreatment()
	app.Run()

	app.Wait()
}
