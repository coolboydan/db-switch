package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"github.com/coolboydan/db-switch/service"
)

var path = flag.String("path", "config.toml", "path")

func main() {

	flag.Parse()

	cfg := service.NewConfig()

	cfg.ConfigFromFile(*path)

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())


	//开启一个新的线程定时检查线程


	//检查是否连接的上,如果连接不上，在3秒内重试3次。

	//检查备用机是否连接的上如果链接的上准备切换，链接不上报错发邮件退出循环。

	//执行切换语句

	//执行完成之后发个邮件


	log.Print("service start")
	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}
