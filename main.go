package main

import (
	"context"

	"smart.gitlab.biomind.com.cn/infrastructure/file-server-gateway/bootstrap"

	_ "go.uber.org/automaxprocs"

	biogo "smart.gitlab.biomind.com.cn/infrastructure/biogo/app"
	"smart.gitlab.biomind.com.cn/infrastructure/biogo/config"
)

func main() {
	ctx := context.Background()

	app := biogo.NewApp(ctx)

	biogo.AddHook(
		config.Init, //初始化配置
		bootstrap.InitLogger,
		bootstrap.InitRedis,
		bootstrap.WatchConn,
		bootstrap.InitRouter(app.HttpEngine), //初始化http router
		bootstrap.InitGrpcServer,             //初始化grpc server
	)

	if err := biogo.Start(); err != nil {
		panic(err)
	}

	err := <-biogo.ErrorCh()
	biogo.Stop()
	panic(err)
}
