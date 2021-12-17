package bootstrap

import (
	"context"
	"file-server-gateway/controller/file"
	"path"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
)

func InitRouter(router *gin.Engine) app.HookFunc {
	return func(ctx context.Context) error {
		dataGrpoup := router.Group("")
		dataGrpoup.GET("/refreshAll", file.RefreshAll)
		dataGrpoup.GET("/refresh/:bucket/:file_name", file.Refresh)
		router.NoRoute(gzip.Gzip(gzip.DefaultCompression), file.Files)

		groupPdf := router.Group("/pdf")
		groupPdf.Static("", path.Join(utils.GetCurrentAbPath(), "data/pdf"))

		//todo 测试用的
		router.NoRoute(gzip.Gzip(gzip.DefaultCompression), file.FilesTest)

		return nil
	}
}
