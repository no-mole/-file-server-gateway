package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"file-server-gateway/model"
	"file-server-gateway/service/dispense"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/output"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/redis"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum/file_server"
	pb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
)

type UrlPath struct {
	Bucket   string `uri:"bucket"`
	FileName string `uri:"file_name"`
}

func Files(ctx *gin.Context) {
	uriPath := ctx.Request.URL.Path
	paths := strings.Split(uriPath, "/")
	if len(paths) < 2 {
		output.Json(ctx, enum.IllegalParam, nil)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)

	p := &UrlPath{
		Bucket:   strings.TrimLeft(strings.Join(paths[:len(paths)-1], "/"), "/"),
		FileName: paths[len(paths)-1],
	}

	//todo 测试用的
	//fileMetadata, err := getFileMetadataFromFile(ctx, p.Bucket, p.FileName)
	//if err != nil {
	//	ctx.Writer.WriteHeader(http.StatusNotFound)
	//	output.Json(ctx, file_server.ErrorGetFileMetadata, err.Error())
	//	return
	//}

	if ctx.Request.Header.Get("Content-Type") == "" {
		ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	} else {
		ctx.Writer.Header().Add("Content-Type", ctx.Request.Header.Get("Content-Type"))
	}

	//todo 测试用的
	//ctx.Writer.Header().Add("e_tage", fileMetadata.ETage)
	//ctx.Writer.Header().Add("header_custom", fileMetadata.Header)
	//ctx.Writer.Header().Add("file_size", fmt.Sprintf("%d", fileMetadata.FileSize))
	//ctx.Writer.Header().Add("file_extension", fileMetadata.FileExtension)

	filePath := path.Join(utils.GetCurrentAbPath(), "data", p.Bucket, p.FileName)
	if exists(filePath) {
		fileOutput(ctx, filePath)
		return
	}

	fileOutputFromNode(ctx, p.Bucket, p.FileName)
}

func fileOutput(ctx *gin.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorFileOpen, err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorFileRead, err.Error())
		return
	}
}

func fileOutputFromNode(ctx *gin.Context, bucket, fileName string) {
	svr := dispense.New()
	download := &pb.DownloadInfo{
		Exist:    false,
		FileName: fileName,
		Bucket:   bucket,
	}
	resp, err := svr.Download(ctx, download)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorDownloadFile, err.Error())
		return
	}

	dirPath := path.Join(utils.GetCurrentAbPath(), "data", bucket)
	if !exists(dirPath) {
		_ = os.MkdirAll(dirPath, os.ModePerm)
	}

	file, err := os.Create(path.Join(dirPath, fileName))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorCreateFile, err.Error())
		return
	}
	defer file.Close()

	_, err = file.Write(resp.Chunk.Content)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorWriteFile, err.Error())
		return
	}

	_, err = io.Copy(ctx.Writer, bytes.NewReader(resp.Chunk.Content))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		output.Json(ctx, file_server.ErrorFileRead, err.Error())
		return
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getFileMetadataFromFile(ctx context.Context, bucket, fileName string) (*pb.UploadFileInfo, error) {
	redis, exist := redis.Client.GetClient(model.RedisEngine)
	if !exist {
		return nil, errors.New("redis not match")
	}
	key := fmt.Sprintf("%s/%s", bucket, fileName)
	body, err := redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	fileMetadata := new(pb.UploadFileInfo)
	err = json.Unmarshal(body, &fileMetadata)
	if err != nil {
		return nil, err
	}
	return fileMetadata, nil
}
