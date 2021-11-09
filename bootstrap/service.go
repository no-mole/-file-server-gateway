package bootstrap

import (
	"context"
	"encoding/json"
	pool "file-server-gateway/grpc_pool"
	"file-server-gateway/model"
	"fmt"
	"math"
	"time"

	fs "smart.gitlab.biomind.com.cn/intelligent-system/biogo/file_server"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/grpc_pool"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config/center"
)

func WatchConn(ctx context.Context) (err error) {
	prefixKey := fmt.Sprintf("%s", model.FileServerNodePrefix)
	fileNodes, err := config.GetClient().GetWithPrefixKey(ctx, prefixKey)
	if err != nil {
		return
	}
	for _, kv := range fileNodes.Kvs {
		s := new(fs.ServerNode)
		err := json.Unmarshal([]byte(kv.Value), &s)
		if err != nil {
			return err
		}
		p, err := pool.NewPool(func() (*grpc.ClientConn, error) {
			return NodeDial(ctx, fmt.Sprintf("%s:%d", s.Host, s.Port))
		})
		pool.ConnMap[kv.Key] = p
		pool.NodeMap[s.NodeName] = s
	}
	_ = pool.LoadLeastNode(ctx, pool.NodeMap)

	config.GetClient().WatchWithPrefix(ctx, fileNodes, func(item *center.Item) {
		str := item.GetValue()
		s := new(fs.ServerNode)
		err := json.Unmarshal([]byte(str), &s)
		if err != nil {
			return
		}

		if item.Act == int64(clientv3.EventTypeDelete) {
			delete(pool.ConnMap, item.Key)
			delete(pool.NodeMap, s.NodeName)
		}
		if item.Act == int64(clientv3.EventTypePut) {

			p, err := pool.NewPool(func() (*grpc.ClientConn, error) {
				return NodeDial(ctx, fmt.Sprintf("%s:%d", s.Host, s.Port))
			})
			if err != nil {
				return
			}
			pool.ConnMap[item.Key] = p
			pool.NodeMap[s.NodeName] = s
		}
		_ = pool.LoadLeastNode(ctx, pool.NodeMap)
	})
	return
}

func NodeDial(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	retryOps := []grpc_retry.CallOption{
		grpc_retry.WithPerRetryTimeout(time.Second * 2),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(time.Second, 0.2)),
	}
	retryInterceptor := grpc_retry.UnaryClientInterceptor(retryOps...)
	StreamRetryInterceptor := grpc_retry.StreamClientInterceptor(retryOps...)
	opts = append([]grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpc_pool.KeepAliveTime,
			Timeout:             grpc_pool.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}, opts...)
	opts = append(opts, grpc.WithChainUnaryInterceptor(retryInterceptor), grpc.WithStreamInterceptor(StreamRetryInterceptor))
	return grpc.DialContext(ctx, target, opts...)
}
