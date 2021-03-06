package grpc_proxy_router

import (
	"fmt"
	"log"
	"net"

	"github.com/LotteWong/giotto-gateway-core/grpc_proxy_middleware"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/reverse_proxy"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

var grpcServers []*grpc.Server

func GrpcServerRun() {
	_, _, grpcServices, err := service.GetSvcService().GroupServicesInMemory()
	if err != nil {
		log.Fatalf(" [ERROR] GrpcServerRun - err:%v\n", err)
	}

	for _, grpcService := range grpcServices {
		tmpService := grpcService
		go func(serviceDetail *po.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GrpcRule.Port)

			lb, err := service.GetLbService().GetLbWithConfForSvc(serviceDetail)
			if err != nil {
				log.Fatalf(" [ERROR] GrpcServerRun - grpc proxy server:%s err:%v\n", addr, err)
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [ERROR] GrpcServerRun - grpc proxy server:%s err:%v\n", addr, err)
			}

			grpcSrvHandler := reverse_proxy.NewGrpcReverseProxy(lb)

			grpcServer := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_proxy_middleware.GrpcJwtAuthMiddleware(serviceDetail),

					grpc_proxy_middleware.GrpcFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcRateLimitMiddleware(serviceDetail),

					grpc_proxy_middleware.GrpcIpAuthMiddleware(serviceDetail),

					grpc_proxy_middleware.GrpcRouteRewriteMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcSrvHandler),
			)
			grpcServers = append(grpcServers, grpcServer)

			log.Printf(" [INFO] GrpcServerRun - grpc proxy server:%s\n", addr)
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf(" [ERROR] GrpcServerRun - grpc proxy server:%s err:%v\n", addr, err)
			}
		}(tmpService)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServers {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] GrpcServerStop - grpc proxy server stopped\n")
	}
}
