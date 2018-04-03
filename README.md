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
  "MaxConcurrencyNum":0, //最大并发处理数 默认值0:进程最大资源打开数量-10

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
      "OpenId": "qwertyuiop",
      "SecretKey": "lqVGFWV42sh5OdxaB-ZMwwRaCVKVcCcQ",
      "BaseUrl": ["http://base.service.soyoung.com/rpc/index"],
      "Type":1
    },
    "frontend": {
      "OpenId": "qwertyuiop",
      "SecretKey": "lqVGFWV42sh5OdxaB-ZMwwRaCVKVcCcQ",
      "BaseUrl": ["http://test.oracle.soyoung.com/rpc/index"],
      "Type":1
    },
    "middleware": {
      "OpenId": "qwertyuiop",
      "SecretKey": "lqVGFWV42sh5OdxaB-ZMwwRaCVKVcCcQ",
      "BaseUrl": ["http://test.admin.oracle.soyoung.com/rpc/index"],
      "Type":1
    }
  },

  "queue":{
    "BrokerList":["172.16.16.4:9092"],
    "Errors":false,
    "Notifications":false
  },

  "topicGroup":{
    "order":{
      "1001":[
        {
          "ServiceName":"order/deposit-order",
          "RpcGroup": "middleware"
        }
      ]
    }
  }
}
```

### Usage
```superconsumer -c /etc/superconsumer.json```