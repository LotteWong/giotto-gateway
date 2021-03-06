package tcp_proxy_middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/service"
	tcp "github.com/LotteWong/tcp-conn-server"
	"github.com/pkg/errors"
)

func TcpRateLimitMiddleware() func(c *tcp.TcpRouterContext) {
	return func(c *tcp.TcpRouterContext) {
		tcpServiceInterface := c.Get("service")
		if tcpServiceInterface == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		tcpServiceDetail := tcpServiceInterface.(*po.ServiceDetail)

		// appInterface := c.Get("app")
		// if appInterface == nil {
		// 	c.Conn.Write([]byte("app not found"))
		// 	c.Abort()
		// 	return
		// }
		// app := appInterface.(*po.App)

		if tcpServiceDetail.AccessControl.ServiceHostFlowLimit != 0 {
			svrServiceName := constants.ServiceFlowCountPrefix + tcpServiceDetail.Info.ServiceName
			svrRateLimit, err := service.GetRateLimitService().GetRateLimit(svrServiceName)
			if err != nil {
				c.Conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			_, _, svrAllow := svrRateLimit.Allow(svrServiceName, tcpServiceDetail.AccessControl.ServiceHostFlowLimit, 1*time.Second)
			// log.Printf("svr name:%s, count:%d\n", svrServiceName, svrCount)
			if !svrAllow {
				c.Conn.Write([]byte(errors.New(fmt.Sprintf("service host flow limit is %d, rate limit exceeds", tcpServiceDetail.AccessControl.ServiceHostFlowLimit)).Error()))
				c.Abort()
				return
			}
		}

		if tcpServiceDetail.AccessControl.ClientIpFlowLimit != 0 {
			var clientIp string
			pair := strings.Split(c.Conn.RemoteAddr().String(), ":")
			if len(pair) != 2 {
				c.Conn.Write([]byte("can not get client ip and port"))
				c.Abort()
				return
			}
			clientIp = pair[0]

			cltServiceName := constants.ServiceFlowCountPrefix + tcpServiceDetail.Info.ServiceName + "_" + clientIp
			cltRateLimit, err := service.GetRateLimitService().GetRateLimit(cltServiceName)
			if err != nil {
				c.Conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			_, _, cltAllow := cltRateLimit.Allow(cltServiceName, tcpServiceDetail.AccessControl.ClientIpFlowLimit, 1*time.Second)
			// log.Printf("clt name:%s, count:%d\n", cltServiceName, cltCount)
			if !cltAllow {
				c.Conn.Write([]byte(errors.New(fmt.Sprintf("client ip flow limit is %d, rate limit exceeds", tcpServiceDetail.AccessControl.ClientIpFlowLimit)).Error()))
				c.Abort()
				return
			}
		}

		// if app.Qps != 0 {
		// 	appServiceName := constants.AppFlowCountPrefix + app.AppId
		// 	appRateLimit, err := service.GetRateLimitService().GetRateLimit(appServiceName)
		// 	if err != nil {
		// 		c.Conn.Write([]byte(err.Error()))
		// 		c.Abort()
		// 		return
		// 	}
		// 	_, _, appAllow := appRateLimit.Allow(appServiceName, app.Qps, 1*time.Second)
		// 	// log.Printf("app name:%s, count:%d\n", appServiceName, appCount)
		// 	if !appAllow {
		// 		c.Conn.Write([]byte(errors.New(fmt.Sprintf("app flow limit is %d, rate limit exceeds", app.Qps)).Error()))
		// 		c.Abort()
		// 		return
		// 	}
		// }

		c.Next()
	}
}
