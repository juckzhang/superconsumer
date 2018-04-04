# superconsumer
super consumer

### feature
```
1、可配置不同消息队列(目前暂支持kafka)。
2、高并发处理消息。
3、可配置最大并发数。
4、监听消费者与业务逻辑处理分离。
5、完善的日志记录系统。
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
    "backend": {
      "OpenId": "OpenId",
      "SecretKey": "SecretKey",
      "BaseUrl": ["http://xxx.com/rpc/index"],
      "Type":1
    },
    "frontend": {
      "OpenId": "OpenId",
      "SecretKey": "SecretKey",
      "BaseUrl": ["http://xxx.com/rpc/index"],
      "Type":1
    },
    "middleware": {
      "OpenId": "OpenId",
      "SecretKey": "SecretKey",
      "BaseUrl": ["http://xxx.com/rpc/index"],
      "Type":1
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
        "ServiceName":"order/deposit-order",
        "RpcGroup": "middleware"
      }
    ]
  }
}
```

### Usage
```superconsumer -c /etc/superconsumer.json```