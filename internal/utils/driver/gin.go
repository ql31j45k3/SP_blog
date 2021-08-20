package driver

import (
	"github.com/ql31j45k3/SP_blog/configs"

	"github.com/gin-gonic/gin"
)

func StartGin(r *gin.Engine) {
	// 控制調試日誌 log
	gin.SetMode(configs.Gin.GetMode())

	r.Run(configs.Host.GetSPBlogAPIHost())
}
