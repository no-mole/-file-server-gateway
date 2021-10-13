package bootstrap

import (
	"context"
	"file-server-gateway/service/dispense"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/middleware"
	fsPb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
)

func InitGrpcServer(_ context.Context) error {
	s := app.NewGrpcServer(
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
