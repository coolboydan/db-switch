package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"github.com/coolboydan/db-switch/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var path = flag.String("path", "config.toml", "path")
var readM *map[string]*sql.DB
var checkTime int

func main() {

	flag.Parse()

	cfg := service.NewConfig()

	cfg.ConfigFromFile(*path)

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())

	//开启一个新的线程定时检查线程
	go doCheck(cfg)
	//检查是否连接的上,如果连接不上，在3秒内重试3次。

	//检查备用机是否连接的上如果链接的上准备切换，链接不上报错发邮件退出循环。

	//执行切换语句

	//执行完成之后发个邮件

	log.Print("service start")
	log.Fatal(http.ListenAndServe("localhost:8881", nil))
}

func doCheck(cfg *service.Config) {

	//持续进行检查

	//检查连续超过3次失败，关闭当前线程启动另一个线程去确认，
	//并且被切换的库也能正常运行
	for true {

		//如果发现连续检查3次失败则进入切换流程
		if checkTime > 3 {
			checkSwitchDb(cfg)
			break
		}

		time.Sleep(2000)
		DBN, err := connect(cfg)

		fmt.Println("name")

		if err != nil {
			fmt.Println(err)
			checkTime++
			continue
		}

		if err := DBN.Ping(); err != nil {
			fmt.Println("open database fail")
			checkTime++
			continue
		}

		checkTime = 0

	}

}

func checkSwitchDb(cfg *service.Config) {
	//只检查一次，如果一次成功则进入切换流程。否则失败
	DBN, err := connect(cfg)

	if err != nil {
		fmt.Println(err)
	}

	if err := DBN.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}

	//执行切换切换流程，1、替换被动配置。2、执行reload语句。3、kill 原有的进程
	err = executeCmd("mv " + cfg.NginxConfig + ".tmp")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = executeCmd("mv " + cfg.SwitchConfig + " " + cfg.NginxConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = executeCmd(cfg.NginxPath + " reload")
	if err != nil {
		fmt.Println(err)
		return
	}

}

func connect(cfg *service.Config) (*sql.DB, error) {
	path := strings.Join([]string{cfg.MonitorDb.User, ":", cfg.MonitorDb.Password, "@tcp(", cfg.MonitorDb.Host, ":", strconv.Itoa(cfg.MonitorDb.Port), ")/", cfg.MonitorDb.Name, "?timeout=5s&readTimeout=6s"}, "")
	DBN, err := connectDB(path)
	return DBN, err
}

func connectDB(path string) (*sql.DB, error) {
	DBN, err := sql.Open("mysql", path)

	DBN.SetConnMaxLifetime(100)

	DBN.SetMaxIdleConns(10)
	return DBN, err
}

func executeCmd(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
