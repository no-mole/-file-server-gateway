package main

import (
	"context"

	"github.com/no-mole/file-server-gateway/bootstrap"

	_ "go.uber.org/automaxprocs"

	biogo "github.com/no-mole/neptune/app"
	"github.com/no-mole/neptune/config"
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
