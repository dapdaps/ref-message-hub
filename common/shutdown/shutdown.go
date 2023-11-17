package shutdown

import (
	"ref-message-hub/common/log"
	"sync/atomic"
	"time"
)

var (
	shutdown       = atomic.Bool{}
	backgroundJobs = make(map[string]*atomic.Bool)
)

func StartJob(name string, workTask func(), cleanupTask func()) {
	go func() {
		registerJob(name)
		for {
			if shutdown.Load() {
				break
			}
			workTask()
		}
		cleanupTask()
		stopJob(name)
	}()
}

func StopAndWaitAll() bool {
	shutdown.Store(true)
	var allFinished bool
	for i := 0; i < 20; i++ {
		allFinished = true
		for name, s := range backgroundJobs {
			if !s.Load() {
				log.Info("background job still not finished: %v", name)
				allFinished = false
			}
		}
		if allFinished {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if allFinished {
		log.Info("all background job finished")
		shutdown.Store(false)
		return true
	} else {
		log.Error("not all tasks shutdown cleanly")
		return false
	}
}

func registerJob(name string) {
	if _, ok := backgroundJobs[name]; ok {
		panic("job already exist: " + name)
	}
	backgroundJobs[name] = &atomic.Bool{}
}

func stopJob(name string) {
	if j := backgroundJobs[name]; j != nil {
		j.Store(true)
		delete(backgroundJobs, name)
		log.Info("job stopped: %v", name)
	} else {
		log.Error("job not exist: %v", name)
	}
}

func StopJob(name string) {
	stopJob(name)
}

//func Shutdown(releaseLockCh chan os.Signal) {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
//	s := <-c
//	if releaseLockCh != nil {
//		common.HandleReleaseLockSignal(releaseLockCh, s)
//	}
//	StopAndWaitAll()
//	log.Info("Service shuts down for signal %v", s)
//}
