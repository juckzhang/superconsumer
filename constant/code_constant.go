package constant

const (
	/** successful operate */
	SUCCESS = 200 //成功状态吗

	/** global process failed */
	UNKNOWN_ERROR             = -1  //未知错误
	PARAM_ERROR               = -2  //参数错误
	REQUEST_METHOD_ERROR      = -3  //请求方式错误
	RPC_PARAM_ERROR           = -4  //rpc调用参数错误
	RPC_APP_FORBIDDEN_A       = -5  //rpcOpenID错误
	RPC_APP_FORBIDDEN_S       = -6  //rpc签名错误
	RPC_CLASS_FORBIDDEN       = -7  //禁止访问
	RPC_CLASS_NOT_EXIST       = -8  //类不存在
	RPC_METHOD_NOT_EXIST      = -9  //方法不存在
	RPC_METHOD_FORBIDDEN      = -10 //禁止访问的
	RPC_FAILED                = -11 //rpc请求失败
	RPC_IP_NOT_ALLOWED        = -12 //rpc禁止访问的ip
	RPC_REQUEST_LOST_EFFICACY = -13 //rpc请求失效
)
