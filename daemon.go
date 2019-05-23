package daemon

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	daemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

var (
	signal = flag.String("s", "", `Send signal to the daemon:
	quit — graceful shutdown
	stop — fast shutdown
	reload — reloading the configuration file
	rotate - rotate the log`)
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
		log.Fatalln("the server is nil")
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

	cntxt := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        args,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	defaultServer = srv

	log.Infoln("daemon started")
	go srv.Serve()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Infoln("daemon terminated")
}
func Status(srv DServer) {
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
	cntxt := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        args,
	}

	d, err := cntxt.Search()
	if err != nil {
		fmt.Println(args, "is stop")
		return
	}
	fmt.Println(args, "[", d.Pid, "]", "is running")
}

func termHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Errorln("nofind server ...")
		return daemon.ErrStop
	}
	if sig == syscall.SIGQUIT {
		log.Debugln("server graceful shutdown ...")
		defaultServer.Quit()
	} else {
		log.Debugln("server fast shutdown ...")
		defaultServer.Stop()
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Errorln("nofind server ...")
		return daemon.ErrStop
	}

	log.Debugln("configuration reloaded")
	defaultServer.Reload()
	return nil
}

func rotateHandler(sig os.Signal) error {
	if defaultServer == nil {
		log.Errorln("nofind server ...")
		return daemon.ErrStop
	}
	log.Debugln("rotate ...")
	defaultServer.Rotate()
	return nil
}
