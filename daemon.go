package daemon

import (
	"flag"
	"log"
	"os"
	"syscall"

	daemon "github.com/sevlyar/go-daemon"
)

var (
	signal = flag.String("s", "", `Send signal to the daemon:
	quit — graceful shutdown
	stop — fast shutdown
	reload — reloading the configuration file
	rotate - rotate the log`)
	foreground    = flag.Bool("f", false, "run at the foreground")
	defaultServer DServer
	defaultArgs   []string
)

const defaultLogFN = "logs/daemon.log"
const defaultPidFN = "run/daemon.pid"

func init() {
	defaultArgs = []string{"[go-daemon sample]"}
}

func Run(srv DServer) {
	if srv == nil {
		log.Panicln("the server is nil")
		return
	}
	flag.Parse()
	pidFileName := srv.GetPidFile()
	if pidFileName == "" {
		pidFileName = defaultPidFN
	}
	logFileName := srv.GetLogFile()
	if logFileName == "" {
		logFileName = defaultLogFN
	}
	args := srv.GetArgs()
	if len(args) == 0 {
		args = defaultArgs
	}
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "rotate"), syscall.SIGUSR1, rotateHandler)

	d := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
		Args:        args,
	}

	if len(daemon.ActiveFlags()) > 0 {
		p, err := d.Search()
		if err != nil {
			log.Panicln("Unable send signal to the daemon: ", err.Error())
		}
		_ = daemon.SendCommands(p)
		return
	}
	if !(*foreground) {
		p, err := d.Reborn()
		if err != nil {
			log.Panicln(err.Error())
		}
		if p != nil {
			return
		}
		defer d.Release()
	}

	defaultServer = srv

	go srv.Serve()

	err := daemon.ServeSignals()
	if err != nil {
		log.Panicln("ServeSignals Error: ", err.Error())
	}

}

func termHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Panicln("the server is nil")
		return daemon.ErrStop
	}
	if sig == syscall.SIGQUIT {
		defaultServer.Quit()
	} else {
		defaultServer.Stop()
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Panicln("the server is nil")
		return daemon.ErrStop
	}

	defaultServer.Reload()
	return nil
}

func rotateHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Panicln("the server is nil")
		return daemon.ErrStop
	}

	defaultServer.Rotate()
	return nil
}
