package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/go-ozzo/ozzo-config"
	"io/ioutil"
	"math/rand"
	"net/http"
	"superconsumer/constant"
	"time"
)

type RpcConfig struct {
	OpenId    string   `服务端的身份id`
	SecretKey string   `请求加密私钥`
	BaseUrl   []string `请求地址`
	Type      int      `请求类型 1: http 2:console脚本`
}

type requestBody struct {
	service   string      `服务名称`
	method    string      `服务方法`
	args      interface{} `请求参数`
	openId    string      `open id`
	timestamp int64       `请求时间`
	sign      string      `请求签名`
}

type request struct {
	data     requestBody
	request  *http.Request
	response *http.Response
}

type RpcClient struct {
	rpcConfig  RpcConfig
	httpClient *http.Client
}

func NewRpcClient(rpcConfig RpcConfig) *RpcClient {
	//校验配置信息是否完整
	if len(rpcConfig.OpenId) == 0 {
		panic("RpcConfig.OpenId is empty!")
	}

	if len(rpcConfig.BaseUrl) == 0 {
		panic("RpcConfig.BaseUrl is empty!")
	}

	if len(rpcConfig.SecretKey) == 0 {
		panic("RpcConfig.SecretKey is empty!")
	}

	if rpcConfig.Type == 0 {
		rpcConfig.Type = 1
	}

	return &RpcClient{
		rpcConfig:  rpcConfig,
		httpClient: &http.Client{},
	}
}

// @description 请求rpc
func (rpcClient *RpcClient) Do(service string, method string, parameters interface{}) int {
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
func (rpcClient *RpcClient) structure(service string, method string, parameters interface{}) (*request, error) {
	//获取配置信息
	data := requestBody{
		service:   service,
		method:    method,
		args:      parameters,
		openId:    rpcClient.rpcConfig.OpenId,
		timestamp: time.Now().Unix(),
		sign:      "asffffasdff",
	}
	request := &request{data: data,}
	dataJson, err := json.Marshal(data)
	var req *http.Request
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(rpcClient.rpcConfig.BaseUrl))
	url := rpcClient.rpcConfig.BaseUrl[index]
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	request.request = req
	return request, nil
}
