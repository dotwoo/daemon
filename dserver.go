package main

// DServer daemon 服务器接口
type DServer interface {
	// Serve 持续性提供服务
	Serve()
	// Quit 优雅关闭服务
	Quit()
	// Stop 快速关闭服务
	Stop()
	// Reload 重载配置
	Reload()
	// Rotate 执行日志 rotate
	Rotate()
	// GetPidFile返回 pid 文件配置
	GetPidFile() string
	// GetLogFile 返回 log 文件配置
	GetLogFile() string
	// TODO 如果父进程的日志设置无法传入子进程,需要
	// SetUplog()
}
