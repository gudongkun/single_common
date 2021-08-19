package common

import (
	"fmt"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func PrometheusBoot(port int) {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		log.Info("metrics listen http server:",fmt.Sprintf(":%d",port))
		err := http.ListenAndServe(fmt.Sprintf(":%d",port), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
