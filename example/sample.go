package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dotwoo/daemon"
	log "github.com/sirupsen/logrus"
)

// SHServer 一个简单的 http 服务 进程配置
type SHServer struct {
	Server    *http.Server
	StartTime string
	PidFN     string
	LogFN     string
	lf        *daemon.FileHandler
}

// Serve 持续性提供服务
func (sh *SHServer) Serve() {
	log.Debugln("shserver start serve ...")
	err := sh.Server.ListenAndServe()
	if err != nil {
		log.Println("Server serve :", err)
	}
	return
}

// Quit 优雅关闭服务
func (sh *SHServer) Quit() {
	log.Debugln("shserver graceful shutdown ...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := sh.Server.Shutdown(ctx)
	if err != nil {
		log.Fatalln("Server quit :", err)
	}
	ctx.Done()
	return
}

// Stop 快速关闭服务
func (sh *SHServer) Stop() {
	log.Debugln("shserver fast stop ...")
	err := sh.Server.Close()
	if err != nil {
		log.Fatalln("Server stop :", err)
	}
	return
}

// Reload 重载配置
func (sh *SHServer) Reload() {
	log.Debugln("shserver reload ...")
	sh.StartTime = time.Now().String()
	return
}

// Rotate 执行日志 rotate
func (sh *SHServer) Rotate() {
	log.Debugln("shserver rotate ...")
	sh.lf.Reopen()
	return
}

// GetPidFile返回 pid 文件配置
func (sh *SHServer) GetPidFile() string {
	// log.Debugln("shserver getPidFile ...")
	return sh.PidFN
}

// GetLogFile 返回 log 文件配置
func (sh *SHServer) GetLogFile() string {
	// log.Debugln("shserver getLogFile ...")
	return sh.LogFN
}

// GetArgs 返回 daemon 参数配置
func (sh *SHServer) GetArgs() []string {
	return []string{"[SHServer sample]"}
}

// NewSample ...
func NewSample() *SHServer {
	sh := new(SHServer)
	sh.PidFN = "./run/shserver.pid"
	sh.LogFN = "./log/shserver.log"
	sh.StartTime = time.Now().String()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page! "+sh.StartTime)
	})
	sh.Server = &http.Server{Addr: ":8566", Handler: mux}
	sh.lf = daemon.NewFileHandler(sh.LogFN, 0640)
	log.SetOutput(sh.lf)
	return sh
}
