package queue

import (
	"github.com/go-ozzo/ozzo-config"
)

type Queue interface{
	Listen(group string)
	Push(topic string, data interface{}) error
}

type Service struct {
	ServiceName string
	GroupRpc string
}

type Topic struct {
	TopicName string
	Services []Service
}

type Group struct {
	GroupId string
	Topics map[string]Topic
}

var (
	groups map[string]Group
	consumer Queue
)

func init() {
	groups = make(map[string]Group)
}

func ConfigGroup(c *config.Config) {

}