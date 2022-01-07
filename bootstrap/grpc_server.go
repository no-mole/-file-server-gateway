package bootstrap

import (
	"context"
	"math"

	"smart.gitlab.biomind.com.cn/infrastructure/file-server-gateway/service/dispense"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/app"
	middleware "smart.gitlab.biomind.com.cn/infrastructure/middlewares"
	fsPb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
)

func InitGrpcServer(_ context.Context) error {
	s := app.NewGrpcServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc_middleware.WithUnaryServerChain(
			middleware.TracingServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			middleware.TracingServerStreamInterceptor(),
		),
	)
	s.RegisterService(&fsPb.Metadata().ServiceDesc, dispense.New())

	return nil
}
