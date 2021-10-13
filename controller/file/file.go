package file

import (
	"file-server-gateway/lru"
	"file-server-gateway/service/dispense"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/output"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum"
	pb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
)

var cache lru.Cache

func init(){
	cache = lru.Constructor(100)
}

type UrlPath struct {
	Bucket string `uri:"bucket"`
	FileName string `uri:"file_name"`
}

func Files(ctx *gin.Context) {
	var urlPath *UrlPath
	if err := ctx.ShouldBindUri(&urlPath);err != nil {
		output.Json(ctx, enum.IllegalParam, err.Error())
		return
	}

	path := cache.Get(path.Join(urlPath.Bucket, urlPath.FileName))
	if path != "" {
		fileOutput(ctx, path.(string))
		return
	}
	fileOutputFromNode(ctx, urlPath.Bucket, urlPath.FileName)
}

func fileOutput(ctx *gin.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
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
	resp, err := svr.Download(ctx,download)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	dirPath := path.Join(utils.GetCurrentAbPath(), "data", bucket)
	if !exists(dirPath) {
		os.MkdirAll(dirPath, os.ModePerm)
	}

	file, err := os.Create(path.Join(dirPath, fileName))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	_, err = file.Write(resp.Chunk.Content)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	cache.Put(path.Join(bucket, fileName),path.Join(dirPath, fileName))
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
