package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"geektime.com/practice/module2/metrics"
	"github.com/golang/glog"
)

func Index(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
		w.Header().Add(k, strings.Join(v, ", "))
	}
	w.Write([]byte("response write ok"))

}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200 ok"))
}

func DelayedHello(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("entering root handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	user := r.URL.Query().Get("user")
	delay := randRange(0, 2000)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}

func randRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
