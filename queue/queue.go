package queue

import (
	"github.com/go-ozzo/ozzo-config"
	"os"
	"os/signal"
)

type Message struct {
	TopicName string
	Data      interface{}
}

type QueueInterface interface {
	Listen()
}

type Queue struct {
	GroupId    string
	Topics     []string
	mChannel   chan Message
	sig        chan os.Signal
	messageNum int `已接收到的消息数量`
}

var (
	sig = make(chan os.Signal, 1)
)

func init() {
	signal.Notify(sig, os.Interrupt) //监听退出信号、user1、user2
}

// @description 创建队列
func NewQueue(c *config.Config, mChan chan Message) QueueInterface {
	var driver string
	if driver = c.GetString("queue.driver"); len(driver) <= 0 {
		panic("message queue driver is empty!")
	}

	var queue QueueInterface
	switch driver {
	case "kafka":
		queue = NewConsumer(c, mChan)
		break
	default:
		break
	}

	return queue
}
