package grpc_pool

import (
	"time"

	"github.com/no-mole/neptune/snowflake"
	"google.golang.org/grpc"
)

type Builder func() (*grpc.ClientConn, error)

var _ ClientConn = (*clientConn)(nil)

type ClientConn interface {
	grpc.ClientConnInterface
	t()
}

type clientConn struct {
	*grpc.ClientConn
}

func (c clientConn) t() {}

var _ Conn = (*rpcConn)(nil)

type Conn interface {
	Conn() ClientConn
	t()
}

type rpcConn struct {
	id      string
	conn    *clientConn
	streams int32
	ts      time.Time
}

func newConn(builder Builder) *rpcConn {
	conn, err := builder()
	if err != nil {
		panic(err)
	}
	return &rpcConn{
		id:   snowflake.GenInt64String(),
		conn: &clientConn{conn},
	}
}

func (conn *rpcConn) Conn() ClientConn {
	return conn.conn
}
func (conn *rpcConn) t() {}
