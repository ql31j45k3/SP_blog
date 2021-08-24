package driver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ql31j45k3/SP_blog/configs"

	"github.com/gin-gonic/gin"
)

func StartGin(cancelCtxStopNotify context.CancelFunc, stopFunc func() context.Context, r *gin.Engine) {
	// 控制調試日誌 log
	gin.SetMode(configs.Gin.GetMode())

	srv := &http.Server{
		Addr:    configs.Host.GetSPBlogAPIHost(),
		Handler: r,
	}

	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"condition": "err != nil and err != http.ErrServerClosed",
				"err":       err,
			}).Error("StartGin - srv.ListenAndServe")
			return
		}
	}(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注意: 業務邏輯有做 goroutine 需用 cancel 通知，確保 goroutine 都有正常中止流程
	cancelCtxStopNotify()

	// 注意: 確認所有開 goroutine 都有完整結束，通知 stop 流程做執行最後關閉資源流程
	ctxStop := stopFunc()
	<-ctxStop.Done()

	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Env.GetShutdownTimeout())
	defer cancelCtx()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("StartGin - srv.Shutdown")
		return
	}

	log.WithFields(log.Fields{
		"shutdownTimeout": fmt.Sprintf("%ds", int64(configs.Env.GetShutdownTimeout()/time.Second)),
	}).Info("StartGin - Server exiting")
}
