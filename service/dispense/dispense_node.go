package dispense

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/no-mole/file-server-gateway/grpc_pool"
	"github.com/no-mole/file-server-gateway/model"

	"github.com/no-mole/neptune/config"

	pb "github.com/no-mole/file-server/protos/file_server"
	"github.com/no-mole/neptune/registry"
)

type Service struct {
	*registry.Metadata
	pb.UnimplementedFileServerServiceServer
}

func New() *Service {
	return &Service{
		Metadata: pb.Metadata(),
	}
}

func (s *Service) SingleUpload(ctx context.Context, in *pb.UploadInfo) (ret *pb.UpLoadResponse, err error) {
	pool := grpc_pool.GetLeastNodePool()
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServerServiceClient(conn.Conn())
	defer pool.Restore(conn)
	return client.SingleUpload(ctx, in)
}

func (s *Service) ChunkUpload(stream pb.FileServerService_ChunkUploadServer) error {
	pool := grpc_pool.GetLeastNodePool()
	conn, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Restore(conn)
	client := pb.NewFileServerServiceClient(conn.Conn())
	putter, err := client.ChunkUpload(context.Background())
	if err != nil {
		return err
	}
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		err = putter.Send(fileChunk)
		if err != nil {
			return err
		}
	}
	resp, err := putter.CloseAndRecv()
	if err != nil {
		return err
	}
	return stream.SendAndClose(resp)
}
func (s *Service) Download(ctx context.Context, in *pb.DownloadInfo) (*pb.DownloadResponse, error) {
	pool, err := grpc_pool.GetNodeConn()
	if err != nil {
		return nil, err
	}
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServerServiceClient(conn.Conn())
	defer pool.Restore(conn)

	resp, err := client.Download(ctx, in)
	if err != nil {
		return nil, err
	}
	if !resp.Exist {
		return s.OtherNodeDownload(ctx, resp, in)
	}
	return resp, nil
}

func (s *Service) OtherNodeDownload(ctx context.Context, resp *pb.DownloadResponse, in *pb.DownloadInfo) (*pb.DownloadResponse, error) {
	key := fmt.Sprintf("/%s/%s/%s", config.GlobalConfig.Namespace,
		model.FileServerNodePrefix, resp.NodeName)
	pool, ok := grpc_pool.ConnMap[key]
	if !ok {
		return nil, errors.New("not match grpc pool")
	}
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServerServiceClient(conn.Conn())
	defer pool.Restore(conn)
	in.Exist = true
	return client.Download(ctx, in)
}

func (s *Service) BigFileDownload(in *pb.DownloadInfo, stream pb.FileServerService_BigFileDownloadServer) error {
	pool, err := grpc_pool.GetNodeConn()
	if err != nil {
		return err
	}
	conn, err := pool.Get()
	if err != nil {
		return err
	}
	client := pb.NewFileServerServiceClient(conn.Conn())
	defer pool.Restore(conn)

	putter, err := client.BigFileDownload(context.Background(), in)
	if err != nil {
		return err
	}
	bigFileDownloadResp := new(pb.DownloadResponse)
	for {
		chunk, err := putter.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !chunk.Exist {
			bigFileDownloadResp = chunk
			break
		}
		stream.Send(&pb.DownloadResponse{
			Chunk:    &pb.Chunk{Content: chunk.Chunk.Content},
			Exist:    true,
			NodeName: "",
		})
	}

	return s.OtherBigFileDownload(in, bigFileDownloadResp.NodeName, stream)
}

func (s *Service) OtherBigFileDownload(in *pb.DownloadInfo, nodeName string, stream pb.FileServerService_BigFileDownloadServer) error {
	key := fmt.Sprintf("/%s/%s/%s", config.GlobalConfig.Namespace,
		model.FileServerNodePrefix, nodeName)
	pool, ok := grpc_pool.ConnMap[key]
	if !ok {
		return errors.New("not match grpc pool")
	}
	conn, err := pool.Get()
	if err != nil {
		return err
	}
	client := pb.NewFileServerServiceClient(conn.Conn())
	defer pool.Restore(conn)
	in.Exist = true
	putter, err := client.BigFileDownload(context.Background(), in)
	if err != nil {
		return err
	}
	for {
		chunk, err := putter.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stream.Send(&pb.DownloadResponse{
			Chunk:    &pb.Chunk{Content: chunk.Chunk.Content},
			Exist:    true,
			NodeName: "",
		})
	}
}
