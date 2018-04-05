# superconsumer
superconsumer 是一个将消息获取与具体业务逻辑处理部分做分离的中间件产品。</br>
核心思想:</br>
通过消息队列获取消息并将消息格式化后放入一个带缓冲的channel。(此处通过channel控制消息处理的并发数)</br>
另一个goroutine从channel中获取消息并根据配置文件中的任务配置做不同处理。比如:通过http调用具体业务处理接口或者调用配置脚本。</br>

解决的问题:</br>
    1、提高消息的并发处理能力。</br>
    2、将具体业务逻辑分离,降低消息处理逻辑的耦合性。</br>
    3、降低业务逻辑串行处理过程中的相互影响.比如:一个处理出错，导致程序退出使另一个逻辑得不到处理。</br>

### Feature
```
1、支持多种消息队列驱动器(目前暂只支持kafka)。
2、高并发处理消息。
3、通过配置管理消息处理的最大并发数。
4、消息监听分发与业务逻辑处理分离。
5、完善的日志记录系统。
6、支持脚本与http接口方式的消息处理接口。
7、实现消息并行化的处理。
```
### Dependency package
```
go get github.com/go-ozzo/ozzo-config
go get github.com/bsm/sarama-cluster
go get github.com/Shopify/sarama
go get github.com/go-ozzo/ozzo-log
```
### Install
```
go get github.com/juckzhang/superconsumer

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
      "OpenId": "OpenId",//业务段分配
      "SecretKey": "SecretKey", //业务端分配
      "BaseUrl": ["http://xxx.com/rpc/index"],
      "Type":1
    },
    "script": {
      "OpenId": "OpenId",//业务段分配
      "SecretKey": "SecretKey",//业务端分配
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

### 业务接口请求参数说明
请求参数中`opendId`,由业务端分配用于做请求校验。
```json
{
    "service":   "服务名称",
    "method":    "服务方法",
    "args":      {
                    "TopicName": "消息topic名称",
                    "Data":      interface{}, //消息内容
                 }
    "openId":    "由服务端分配的openId"
    "timestamp": "请求发起时间"
    "sign":      "请求签名"
}
```

### Usage
```
启动服务: superconsumerd start
停止服务: superconsumerd stop
重启服务: superconsumerd restart
重新加载配置: superconsumerd reload
```