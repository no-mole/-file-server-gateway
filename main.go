package main

import (
	"context"

	_ "go.uber.org/automaxprocs"

	"file-server-gateway/bootstrap"

	biogo "smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
)

func main() {
	ctx := context.Background()

	app := biogo.NewApp(ctx)

	biogo.AddHook(
		config.Init, //初始化配置
		bootstrap.InitLogger,
		bootstrap.WatchConn,
		bootstrap.InitRouter(app.HttpEngine), //初始化http router
		bootstrap.InitGrpcServer, //初始化grpc server
	)

	if err := biogo.Start(); err != nil {
		panic(err)
	}

	err := <-biogo.ErrorCh()
	biogo.Stop()
	panic(err)
}
