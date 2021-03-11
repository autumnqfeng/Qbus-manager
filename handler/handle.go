package handler

import (
	"Qbus-manager/pkg/errno"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ErrCode    	string           `json:"errcode"`
	ErrMsg 		string      	 `json:"errmsg"`
	Result    	interface{}      `json:"result"`
}

type ResponseV1 struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, result interface{})  {
	code, message := errno.DecodeErr(err)

	errCode := "Failed"
	if code == errno.OK.Code {
		errCode = "OK"
		message = ""
	}

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		ErrCode:    errCode,
		ErrMsg: 	message,
		Result:    	result,
	})
}

func SendResponseV1(c *gin.Context, err error, result interface{})  {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, ResponseV1{
		Code:    	code,
		Message: 	message,
		Data:    	result,
	})
}