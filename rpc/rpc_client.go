package rpc

import (
	"github.com/go-ozzo/ozzo-config"
	"net/http"
	"bytes"
	"encoding/json"
	"time"
	"math/rand"
	"superconsume/constant"
	"io/ioutil"
	"fmt"
	"os"
)

type RpcConfig struct{
	OpenId string `服务端的身份id`
	SecretKey string `请求加密私钥`
	BaseUrl []string `请求地址`
	Type int `请求类型 1: http 2:console脚本`
}

type requestBody struct {
	service string `服务名称`
	method  string `服务方法`
	args interface{} `请求参数`
	openId string `open id`
	timestamp int64 `请求时间`
	sign string `请求签名`
}

type request struct {
	data    requestBody
	request *http.Request
	response *http.Response
}

type RpcClient struct {
	rpc_config RpcConfig
	httpClient *http.Client
}

var (
	rpc_config_group map[string]RpcConfig
	rpc_client map[string]*RpcClient
)

func Config(c *config.Config)  {
	c.Configure(&rpc_config_group, "rpc")
	fmt.Print(rpc_config_group);os.Exit(0);
	if len(rpc_config_group) == 0 {
		panic("rpc配置信息不合法")
	}

	rpc_client = make(map[string]*RpcClient)
	for key, val := range rpc_config_group {

		rpc_client[key] = &RpcClient{
			rpc_config:val,
			httpClient:&http.Client{},
		}
	}
}

// @description 创建rpc客户端
func NewRpcClient(group string) (*RpcClient,error){
	if rpcClient, ok := rpc_client[group]; ok {
		return rpcClient, nil
	}

	panic("rpc client 不存在!")
}

// @description 请求rpc
func (rpcClient *RpcClient) Do(service string, method string, parameters interface{}) (int) {
	var req *request
	var err error
	if req, err = rpcClient.structure(service, method, parameters); err != nil {
		return constant.RPC_FAILED
	}
	req.response, err = rpcClient.httpClient.Do(req.request)
	if err != nil {
		return constant.RPC_FAILED
	}

	if req.response.StatusCode != 200 {
		return constant.RPC_FAILED
	}
	defer req.response.Body.Close()
	body, _ := ioutil.ReadAll(req.response.Body)
	c := config.New()
	c.LoadJSON(body)

	return c.GetInt("code")
}

// @description 构造请求
func (rpcClient *RpcClient) structure(service string, method string, parameters interface{}) (*request,error) {
	//获取配置信息
	data := requestBody{
		service: service,
		method: method,
		args:parameters,
		openId:rpcClient.rpc_config.OpenId,
		timestamp:time.Now().Unix(),
		sign:"asffffasdff",
	}
	request := &request{
		data:data,
		request:new(http.Request),
		response:new(http.Response),
	}
	dataJson, err := json.Marshal(data)
	var req *http.Request
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(rpcClient.rpc_config.BaseUrl))
	url := rpcClient.rpc_config.BaseUrl[index]
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil,err
	}

	request.request = req
	return request, nil
}