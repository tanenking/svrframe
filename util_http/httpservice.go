package util_http

import (
	"fmt"
	"net/http"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/logx"

	"github.com/gin-gonic/gin"
)

func update() {
	go func() {
		exit_ch := constants.GetServiceStopListener().AddListener()
		<-exit_ch.Done()
		httpServer.Close()
		httpServer = nil
	}()
}

func StartHttpService(g *gin.Engine) {

	serviceConfig := config.GetServiceConfig()
	port := serviceConfig.Service.HttpPort
	if port == 0 {
		port = 80
	}
	addr := fmt.Sprintf(":%d", port)
	httpServer = &http.Server{
		Addr:    addr,
		Handler: g,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.ErrorF("%v", err)
		}
	}()

	logx.InfoF("HTTP SERVER [ %s ] RUNNING", constants.Service_Type)

	update()
}
