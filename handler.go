package main

import (
	"github.com/kataras/iris/v12"
	"net"
	"net/http"
	"strings"
)

func CounterHandler(ctx iris.Context) {
	//println(ctx.Path())//不包括域名
	//println(ctx.RequestPath(true)) //同上
	//println(ctx.RequestPath(false)) //同上
	//println(ctx.FullRequestURI()) // 包括域名

	//访问根路径时统计ip和访问量
	if strings.Compare(ctx.Path(), "/") == 0 {
		var i float64
		i = 1

		//统计ip
		if _, ok := (ApplicationCache["visitors"].(map[string]interface{}))[clientIP(ctx.Request())]; ok {
			i = (ApplicationCache["visitors"].(map[string]interface{}))[clientIP(ctx.Request())].(float64) + 1

		}
		(ApplicationCache["visitors"].(map[string]interface{}))[clientIP(ctx.Request())] = i
		//增加访问量
		i = ApplicationCache["visitorCount"].(float64)
		ApplicationCache["visitorCount"] = i + 1
		UpdateApplicationJson()
	}
	ctx.Next()
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func clientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
