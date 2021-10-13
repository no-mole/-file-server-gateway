package file

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/output"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum"
	"smart.gitlab.biomind.com.cn/intelligent-system/enum/file_server"
)

func Refresh(ctx *gin.Context) {
	var urlPath *UrlPath
	if err := ctx.ShouldBindUri(&urlPath);err != nil {
		output.Json(ctx, enum.IllegalParam, err.Error())
		return
	}

	cache.Remove(path.Join(urlPath.Bucket, urlPath.FileName))
	err := os.RemoveAll(path.Join(utils.GetCurrentAbPath(),"data", urlPath.Bucket, urlPath.FileName))
	if err != nil {
		output.Json(ctx, file_server.ErrorRemoveFile, nil)
		return
	}
	output.Json(ctx, enum.Success, nil)
}

func RefreshAll(ctx *gin.Context) {
	dir ,err := ioutil.ReadDir(path.Join(utils.GetCurrentAbPath(),"data"))
	if err != nil {
		output.Json(ctx, file_server.ErrorDirOpen, err.Error())
		return
	}
	for _, d := range dir {
		os.RemoveAll(path.Join(utils.GetCurrentAbPath(),"data", d.Name()))
	}
	cache = cache.Refresh()
	output.Json(ctx, enum.Success, nil)
}
