package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang/glog"

	"geektime.com/practice/module2/handler"
	"geektime.com/practice/module2/middleware"
)

var portFlag int
var pidFlag int

func init() {
	flag.IntVar(&portFlag, "port", 3000, "port number")
	flag.IntVar(&pidFlag, "pid", 0, "pid number")
}

func gracefullyExit(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	glog.Info("server will shutdown in 5 seconds")
	// https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
	if err := server.Shutdown(ctx); err != nil {
		glog.Fatal(err)
	}
}

// https://cloud.tencent.com/developer/article/1645996
func dealSysSignal(server *http.Server, wg *sync.WaitGroup, c chan os.Signal) {
Loop:
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			break Loop
		}
	}
	gracefullyExit(server)
	wg.Done()
}

func newServerHandler() *http.ServeMux {
	interceptor := func(hf http.HandlerFunc) http.HandlerFunc {
		return middleware.SetVersion(middleware.GetIP(middleware.GetResponseStatus(hf)))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", interceptor(handler.Index))
	mux.HandleFunc("/healthz", interceptor(handler.Healthz))
	return mux
}

func main() {
	flag.Parse()

	wg := &sync.WaitGroup{}
	server := http.Server{Addr: fmt.Sprintf(":%d", portFlag), Handler: newServerHandler()}

	wg.Add(1)
	server.RegisterOnShutdown(func() {
		wg.Done()
		glog.Info("server is shutdown")
	})

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go dealSysSignal(&server, wg, c)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		glog.Fatal(err)
	}

	wg.Wait()
	// avoid send close channel
	signal.Stop(c)
	close(c)
	os.Exit(0)
}
