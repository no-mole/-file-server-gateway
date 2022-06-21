package bootstrap

import (
	"context"
	"path"

	"github.com/no-mole/file-server-gateway/controller/file"

	"github.com/no-mole/neptune/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/no-mole/neptune/app"
)

func InitRouter(router *gin.Engine) app.HookFunc {
	return func(ctx context.Context) error {
		dataGrpoup := router.Group("")
		dataGrpoup.GET("/refreshAll", file.RefreshAll)
		dataGrpoup.GET("/refresh/:bucket/:file_name", file.Refresh)
		router.NoRoute(gzip.Gzip(gzip.DefaultCompression), file.Files)

		groupPdf := router.Group("/pdf")
		groupPdf.Static("", path.Join(utils.GetCurrentAbPath(), "data/pdf"))

		return nil
	}
}
