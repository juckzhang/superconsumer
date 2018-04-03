package queue

import (
	"github.com/go-ozzo/ozzo-config"
	"github.com/bsm/sarama-cluster"
	"github.com/Shopify/sarama"
	"superconsume/log"
	"sync"
	"superconsume/rpc"
)

type kafkaConfig struct {
	BrokerList []string `片列表`
	Errors bool
	Notifications bool
}

type Kafka struct {
	kafka_config kafkaConfig
	groupId string
	topic map[string][]Service
	consumer *cluster.Consumer
	c *config.Config
}

func NewKafkaConsumer(c *config.Config) *Kafka {
	var kafka_config kafkaConfig
	err := c.Configure(&kafka_config, "queue.driver_config")
	if err != nil {
		panic(err)
	}

	return &Kafka{
		kafka_config:kafka_config,
		c: c,
		topic:make(map[string][]Service),
	}
}
func (kafka *Kafka) Listen(group string) {
	defer wg.Done()
	var err error
	//根据groupId获取相应的配置信息
	kafka.groupId = group
	if err = kafka.c.Configure(&kafka.topic, "topicGroup."+group); err != nil {
		panic(err)
	}

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = kafka.kafka_config.Errors
	config.Group.Return.Notifications = kafka.kafka_config.Notifications
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	topics := make([]string, 1)
	for key,_ := range kafka.topic {
		topics = append(topics, key)
	}
	// init consumer
	kafka.consumer, err = cluster.NewConsumer(kafka.kafka_config.BrokerList, group, topics, config)
	if err != nil {
		log.Warning("application","%s: sarama.NewSyncProducer err, message=%s \n", group, err)
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
				log.Warning("application", "%s:%s/%d/%d\t%s\t%s\n", group, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				kafka.consumer.MarkOffset(msg, "")  // mark message as processed
				//处理消息
				for task := range kafka.topic[msg.Topic] {
					if rpc_client, err:= rpc.NewRpcClient(task.RpcGroup); err != nil {
						panic(err)

						go rpc_client.Do(task)
					}
				}
			}
		case <-signals:
			break Loop
		}
	}

	sw.Wait()
}

func (kafka *Kafka) errors()  {
	// consume errors
	go func() {
		for err := range kafka.consumer.Errors() {
			log.Warning("application","%s:Error: %s\n", kafka.groupId, err.Error())
		}
	}()
}

func (kafka *Kafka) notification() {
	// consume notifications
	go func() {
		for ntf := range kafka.consumer.Notifications() {
			log.Warning("application","%s:Rebalanced: %+v \n", kafka.groupId, ntf)
		}
	}()
}
