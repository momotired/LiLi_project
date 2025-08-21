package middleware

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// CORS中间件
func CORS(ctx *context.Context) {
	// 获取CORS配置
	allowOrigins, _ := beego.AppConfig.String("cors_allow_origins")
	allowMethods, _ := beego.AppConfig.String("cors_allow_methods")
	allowHeaders, _ := beego.AppConfig.String("cors_allow_headers")

	// 设置默认值
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	if allowMethods == "" {
		allowMethods = "GET,POST,PUT,DELETE,OPTIONS"
	}
	if allowHeaders == "" {
		allowHeaders = "Content-Type,Authorization,X-Requested-With"
	}

	// 获取请求的Origin
	origin := ctx.Request.Header.Get("Origin")

	// 如果配置了具体的域名，检查Origin是否在允许列表中
	if allowOrigins != "*" && origin != "" {
		origins := strings.Split(allowOrigins, ",")
		originAllowed := false
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				originAllowed = true
				break
			}
		}
		if originAllowed {
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
		}
	} else {
		ctx.Output.Header("Access-Control-Allow-Origin", allowOrigins)
	}

	// 设置其他CORS头
	ctx.Output.Header("Access-Control-Allow-Methods", allowMethods)
	ctx.Output.Header("Access-Control-Allow-Headers", allowHeaders)
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Max-Age", "86400") // 缓存预检请求结果24小时

	// 处理预检请求
	if ctx.Request.Method == "OPTIONS" {
		ctx.Output.SetStatus(200)
		ctx.Output.Body([]byte(""))
		return
	}
}

// 全局CORS中间件
func GlobalCORS(ctx *context.Context) {
	// 对所有请求应用CORS
	CORS(ctx)
}

// 访问日志中间件
func AccessLog(ctx *context.Context) {
	// 记录请求开始时间
	startTime := time.Now()
	
	// 获取请求信息
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	clientIP := ctx.Input.IP()
	userAgent := ctx.Request.Header.Get("User-Agent")
	
	// 记录访问日志
	logs.Info("[ACCESS] %s %s - IP: %s - UA: %s - Time: %s", 
		method, path, clientIP, userAgent, startTime.Format("2006-01-02 15:04:05"))
}
