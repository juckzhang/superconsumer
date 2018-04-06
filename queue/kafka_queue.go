package queue

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/go-ozzo/ozzo-config"
	"superconsumer/log"
)

type Kafka struct {
	Queue
	BrokerList    []string
	Errors        bool
	Notifications bool
	consumer      *cluster.Consumer
}

func NewConsumer(c *config.Config, ch chan Message) *Kafka {
	var kafkaClient *Kafka
	err := c.Configure(&kafkaClient, "queue.driver_config")
	if err != nil {
		panic(err)
	}
	kafkaClient.mChannel = ch
	kafkaClient.sig = sig
	return kafkaClient
}

func (kafka *Kafka) Listen() {
	var err error

	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = kafka.Errors
	cfg.Group.Return.Notifications = kafka.Notifications
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	if kafka.consumer, err = cluster.NewConsumer(kafka.BrokerList, kafka.GroupId, kafka.Topics, cfg); err != nil {
		panic(err)
	}
	defer kafka.consumer.Close()

	kafka.errors()
	kafka.notification()

Loop:
	for {
		select {
		case msg, ok := <-kafka.consumer.Messages():
		    if !ok {break Loop}
            log.Debug("application", "%s:%s/%d/%d\t%s\t%s\n", kafka.GroupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
            kafka.consumer.MarkOffset(msg, "") // mark message as processed
            kafka.messageNum++                 //接收到的消息数量
            c := config.New()
            c.LoadJSON(msg.Value)
            message := Message{
                TopicName: msg.Topic,
                Data:      c.Get("data"),
            }
            kafka.mChannel <- message //将消息如管道，同时管道具有控制消息并发处理数的作用
		case <-kafka.sig:
			break Loop
		}
	}
    close(kafka.mChannel)
}

func (kafka *Kafka) errors() {
	// consume errors
	go func() {
		for err := range kafka.consumer.Errors() {
			log.Warning("application", "%s:Error: %s\n", kafka.GroupId, err.Error())
		}
	}()
}

func (kafka *Kafka) notification() {
	// consume notifications
	go func() {
		for ntf := range kafka.consumer.Notifications() {
			log.Warning("application", "%s:Rebalanced: %+v \n", kafka.GroupId, ntf)
		}
	}()
}
