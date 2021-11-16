package bootstrap

import (
	"context"
	"file-server-gateway/controller/file"
	"github.com/gin-gonic/gin"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/app"
)


func InitRouter(router *gin.Engine) app.HookFunc {
	return func(ctx context.Context) error {
		dataGrpoup := router.Group("")
		dataGrpoup.GET("/:bucket/:file_name", file.Files)
		dataGrpoup.GET("/refreshAll", file.RefreshAll)
		dataGrpoup.GET("/refresh/:bucket/:file_name", file.Refresh)

		return nil
	}
}