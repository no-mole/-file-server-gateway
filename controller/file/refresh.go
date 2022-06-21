package file

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	fileserver "github.com/no-mole/file-server-gateway/enum"
	"github.com/no-mole/neptune/enum"
	"github.com/no-mole/neptune/output"
	"github.com/no-mole/neptune/utils"
)

func Refresh(ctx *gin.Context) {
	var urlPath *UrlPath
	if err := ctx.ShouldBindUri(&urlPath); err != nil {
		output.Json(ctx, enum.IllegalParam, err.Error())
		return
	}

	err := os.RemoveAll(path.Join(utils.GetCurrentAbPath(), "data", urlPath.Bucket, urlPath.FileName))
	if err != nil {
		output.Json(ctx, fileserver.ErrorRemoveFile, nil)
		return
	}
	output.Json(ctx, enum.Success, nil)
}

func RefreshAll(ctx *gin.Context) {
	dir, err := ioutil.ReadDir(path.Join(utils.GetCurrentAbPath(), "data"))
	if err != nil {
		output.Json(ctx, fileserver.ErrorDirOpen, err.Error())
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(path.Join(utils.GetCurrentAbPath(), "data", d.Name()))
	}
	output.Json(ctx, enum.Success, nil)
}
