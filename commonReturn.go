package Hornetgo

//CommonReturn 返回数据格式
type CommonReturn struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

//SendError 返回错误
func SendError(code int, msg string) CommonReturn {

	return CommonReturn{
		ErrCode: code,
		ErrMsg:  msg,
	}

}

//SendSuccess 操作成功
func SendSuccess() CommonReturn {

	return CommonReturn{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    "success",
	}

}

//SendResult 返回数据
func SendResult(v interface{}) CommonReturn {

	return CommonReturn{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    v,
	}

}
