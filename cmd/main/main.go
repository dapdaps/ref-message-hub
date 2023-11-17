package main

import (
	"flag"
	"os"
	"os/signal"
	"ref-message-hub/common/log"
	"ref-message-hub/common/shutdown"
	"ref-message-hub/conf"
	"ref-message-hub/http"
	"ref-message-hub/service"
	"syscall"
)

// nolint:all
func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Init(conf.Conf.Log, conf.Conf.Debug)
	log.Info("message-hub service start")

	service.Init()
	log.Info("message-hub service init")

	http.InitHttpServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	s := <-c
	log.Info("message-hub exit for signal %v", s)

	http.ShutdownHttpServer()
	shutdown.StopAndWaitAll()
}
