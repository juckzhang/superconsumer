package queue

import (
	"github.com/go-ozzo/ozzo-config"
	"os"
	"os/signal"
	"sync"
)

type Queue interface{
	Listen(group string)
}

type Service struct {
	ServiceName string
	GroupRpc string
}

var (
	groups map[string]Queue
	signals chan os.Signal
	wg sync.WaitGroup
)

func init() {
	groups = make(map[string]Queue)
	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
}

func Listen(c *config.Config) {
	for groupId, _ := range c.Get("topicGroup").(map[string]interface{}) {
		queue := newQueue(c)
		groups[groupId] = queue
		wg.Add(1)
		go queue.Listen(groupId)
	}
}

// @description 创建队列
func newQueue(c *config.Config) Queue{
	var driver string
	if driver = c.GetString("queue.driver"); len(driver) <= 0 {
		panic(driver)
	}

	var queue Queue
	switch driver {
	case "kafka":
		queue =  NewKafkaConsumer(c)
		break
	default:
		break
	}

	return queue
}