package bootstrap

import (
	"context"
	"math"

	"github.com/no-mole/file-server-gateway/service/dispense"

	"google.golang.org/grpc"

	fsPb "github.com/no-mole/file-server/protos/file_server"
	"github.com/no-mole/neptune/app"
)

func InitGrpcServer(_ context.Context) error {
	s := app.NewGrpcServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
	)
	s.RegisterService(&fsPb.Metadata().ServiceDesc, dispense.New())

	return nil
}
