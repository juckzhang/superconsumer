package main

import (
	"flag"
	"github.com/go-ozzo/ozzo-config"
	"superconsume/rpc"
	"superconsume/log"
	"superconsume/queue"
)

type App struct {
	C *config.Config
	Config_path *string
}
var (
	app = App{
		C: config.New(),
		Config_path:flag.String("c", "/etc/superconsumer.json", "配置文件"),
	}
)

func main() {
	flag.Parse()

	//加载配置文件
	if err := app.C.Load(*app.Config_path); err != nil {
		panic(err)
	}
	//配置logger
	log.NewLogger(app.C)
	//配置rpc
	rpc.Config(app.C)

	//加载队列配置
	queue.Listen(app.C)
}