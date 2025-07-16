package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"` // 链路追踪ID
	Details   interface{} `json:"details,omitempty"`  // 错误详情
}

func (r Response) Error() string {
	return fmt.Sprintf("code: %d, message: %s", r.Code, r.Message)
}

// 基础状态码常量
const (
	SUCCESS         = 200
	ERROR_PARAM     = 400
	ERROR_AUTH      = 401
	ERROR_FORBID    = 403
	ERROR_NOT_FOUND = 404
	ERROR_SERVER    = 500
	ERROR_WECHAT    = 1001
)

// 业务错误码常量
const (
	ERROR_USER_NOT_FOUND      = 2001 // 用户不存在
	ERROR_USER_DISABLED       = 2002 // 用户被禁用
	ERROR_USER_ALREADY_EXISTS = 2003 // 用户已存在
	ERROR_TOKEN_EXPIRED       = 2004 // Token过期
	ERROR_TOKEN_INVALID       = 2005 // Token无效
	ERROR_CODE_INVALID        = 2006 // 微信Code无效
	ERROR_DECRYPT_FAILED      = 2007 // 解密失败
	ERROR_DATABASE            = 2008 // 数据库错误
	ERROR_EXTERNAL_API        = 2009 // 外部API调用失败
)

// 错误信息映射
var ErrorMessages = map[int]string{
	SUCCESS:                   "success",
	ERROR_PARAM:               "参数错误",
	ERROR_AUTH:                "认证失败",
	ERROR_FORBID:              "权限不足",
	ERROR_NOT_FOUND:           "资源不存在",
	ERROR_SERVER:              "服务器内部错误",
	ERROR_WECHAT:              "微信接口错误",
	ERROR_USER_NOT_FOUND:      "用户不存在",
	ERROR_USER_DISABLED:       "用户已被禁用",
	ERROR_USER_ALREADY_EXISTS: "用户已存在",
	ERROR_TOKEN_EXPIRED:       "Token已过期",
	ERROR_TOKEN_INVALID:       "Token无效",
	ERROR_CODE_INVALID:        "微信授权码无效",
	ERROR_DECRYPT_FAILED:      "数据解密失败",
	ERROR_DATABASE:            "数据库操作失败",
	ERROR_EXTERNAL_API:        "外部接口调用失败",
}

// 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:      SUCCESS,
		Message:   ErrorMessages[SUCCESS],
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// 错误响应
func Error(code int, message string) Response {
	return Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}

// 错误响应（使用预定义消息）
func ErrorWithCode(code int) Response {
	message := ErrorMessages[code]
	if message == "" {
		message = "未知错误"
	}
	return Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}

// 带详情的错误响应
func ErrorWithDetails(code int, message string, details interface{}) Response {
	return Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Details:   details,
		Timestamp: time.Now().Unix(),
	}
}

// 输出JSON响应
func WriteJSON(ctx *context.Context, resp Response) {
	ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	data, _ := json.Marshal(resp)
	ctx.Output.Body(data)
}

// 输出成功响应
func WriteSuccess(ctx *context.Context, data interface{}) {
	WriteJSON(ctx, Success(data))
}

// 输出错误响应
func WriteError(ctx *context.Context, code int, message string) {
	// 记录错误日志
	logs.Error("API Error: code=%d, message=%s, uri=%s", code, message, ctx.Request.RequestURI)
	WriteJSON(ctx, Error(code, message))
}

// 输出错误响应（使用预定义消息）
func WriteErrorWithCode(ctx *context.Context, code int) {
	message := ErrorMessages[code]
	if message == "" {
		message = "未知错误"
	}
	logs.Error("API Error: code=%d, message=%s, uri=%s", code, message, ctx.Request.RequestURI)
	WriteJSON(ctx, ErrorWithCode(code))
}

// 输出带详情的错误响应
func WriteErrorWithDetails(ctx *context.Context, code int, message string, details interface{}) {
	logs.Error("API Error: code=%d, message=%s, details=%+v, uri=%s", code, message, details, ctx.Request.RequestURI)
	WriteJSON(ctx, ErrorWithDetails(code, message, details))
}

// 业务异常结构
type BusinessError struct {
	Code    int
	Message string
	Details interface{}
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("Business Error: %d - %s", e.Code, e.Message)
}

// 创建业务异常
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// 创建带详情的业务异常
func NewBusinessErrorWithDetails(code int, message string, details interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// 处理业务异常
func HandleBusinessError(ctx *context.Context, err error) {
	if bizErr, ok := err.(*BusinessError); ok {
		if bizErr.Details != nil {
			WriteErrorWithDetails(ctx, bizErr.Code, bizErr.Message, bizErr.Details)
		} else {
			WriteError(ctx, bizErr.Code, bizErr.Message)
		}
	} else {
		WriteError(ctx, ERROR_SERVER, err.Error())
	}
}
