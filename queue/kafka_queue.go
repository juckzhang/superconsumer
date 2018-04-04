package queue

import (
    "github.com/Shopify/sarama"
    "github.com/bsm/sarama-cluster"
    "github.com/go-ozzo/ozzo-config"
    "superconsume/log"
    "sync"
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
    return kafkaClient
}

func (kafka *Kafka) Listen() {
    var sw sync.WaitGroup
    var err error

    cfg := cluster.NewConfig()
    cfg.Consumer.Return.Errors = kafka.Errors
    cfg.Group.Return.Notifications = kafka.Notifications
    cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

    kafka.consumer, err = cluster.NewConsumer(kafka.BrokerList, kafka.GroupId, kafka.Topics, cfg)
    if err != nil {
        log.Warning("application", "%s: sarama.NewSyncProducer err, message=%s \n", kafka.GroupId, err)
        return
    }
    defer kafka.consumer.Close()

    kafka.errors()
    kafka.notification()

Loop:
    for {
        select {
        case msg, ok := <-kafka.consumer.Messages():
            if ok {
                message := Message{
                    TopicName: msg.Topic,
                    Data:      msg.Value,
                }
                log.Info("application", "%s:%s/%d/%d\t%s\t%s\n", kafka.GroupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
                kafka.consumer.MarkOffset(msg, "") // mark message as processed
                kafka.messageNum++                 //接收到的消息数量
                kafka.mChannel <- message          //将消息如管道，同时管道具有控制消息并发处理数的作用
            }
        case <-kafka.sig:
            break Loop
        }
    }

    sw.Wait()
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
