# superconsumer
superconsumer 是一个将消息队列消费者与具体业务逻辑做分离的中间件产品。</br>
其核心思想是:通过消息队列获取消息并将消息格式化后放入一个带缓冲的channel。(此处通过channel控制消息处理的并发数)</br>
              另一个goroutine从channel中获取消息并根据配置文件中的任务配置做不同处理。比如:通过http调用具体业务处理接口或者调用配置脚本。</br>
解决的问题:</br>
    1、提高消息的并发处理能力。</br>
    2、将具体业务逻辑分离,降低消息处理逻辑的耦合性。</br>
    3、降低业务逻辑串行处理过程中的相互影响.比如:一个处理出错，导致程序退出使另一个逻辑得不到处理。</br>

### Feature
```
1、可配置不同消息队列(目前暂只支持kafka)。
2、高并发处理消息。
3、可配置最大并发数。
4、监听消费者与业务逻辑处理分离。
5、完善的日志记录系统。
6、支持脚本与http接口方式的任务处理接口
```
### Dependency package
```
go get github.com/go-ozzo/ozzo-config
go get github.com/bsm/sarama-cluster
go get github.com/Shopify/sarama
go get github.com/go-ozzo/ozzo-log
```

### Config
```json
{
  "maxConcurrent":1000, //最大并发处理数 默认值0:进程最大资源打开数量-10

  //日志组件
  "logger": {
    "Targets": [
      {
        "type":     "FileTarget",
        "FileName": "runtime/log/app.log",
        "MaxLevel": 4 // Warning or above
      }
    ]
  },

  "rpc":{
    "http": {
      "OpenId": "OpenId",
      "SecretKey": "SecretKey",
      "BaseUrl": ["http://xxx.com/rpc/index"],
      "Type":1
    },
    "script": {
      "OpenId": "OpenId",
      "SecretKey": "SecretKey",
      "BaseUrl": ["php -f script.php"],
      "Type":2
    }
  },

  "queue":{
    "driver":"kafka",
    "driver_config":{
      "Group":"oracle-kafka-queue-all",
      "BrokerList":["127.0.0.1:9092"],
      "Errors":false,
      "Notifications":false,
      "Topics":["1001"]
    }
  },

  "topicGroup":{
    "1001":[
      {
        "ServiceName":"service",
        "MethodName": "method",
        "RpcGroup": "http" //对应与`rpc`配置下的具体客户端
      },
      {
        "ServiceName":"service",
        "MethodName": "method",
        "RpcGroup": "script" //对应与`rpc`配置下的具体客户端
      }
    ]
  }
}
```

### Usage
```
启动服务: superconsumer start
停止服务: superconsumer stop
重启服务: superconsumer restart
重新加载配置: superconsumer reload
```