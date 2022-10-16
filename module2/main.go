package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

func main1() {

	que := make(chan int, 10)
	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	producerNum := 10
	customerNum := 10
	wg.Add(producerNum + customerNum)
	// producer
	for i := 0; i < producerNum; i++ {
		go func(n int) {
			for {
				select {
				case <-ticker.C:
					que <- n
				case <-ctx.Done():
					wg.Done()
					fmt.Printf("producer %d exit\n", n)
					return
				}
			}

		}(i)
	}

	// consumer
	for i := 0; i < customerNum; i++ {
		go func(n int) {
			for {
				select {
				case v := <-que:
					fmt.Printf("consume data %d\n", v)
				case <-ctx.Done():
					wg.Done()
					fmt.Printf("consumer %d exit\n", n)
					return
				}
			}
		}(i)

	}

	time.Sleep(20 * time.Second)
	cancel()
	wg.Wait()

}

// 参考：https://gist.github.com/Boerworz/b683e46ae0761056a636
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// getIP returns the ip address from the http request
// 参考：https://gist.github.com/miguelmota/7b765edff00dc676215d6174f3f30216
func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

func interceptor(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		w.Header().Add("version", os.Getenv("VERSION"))

		lrw := NewLoggingResponseWriter(w)

		handler(lrw, r)

		// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
		addr, err := getIP(r)
		if err != nil {
			glog.Errorf("get addr error %v\n", err)
		}
		glog.Infof("response status code %d, ip addr %s\n", lrw.statusCode, addr)

	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", interceptor(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
			w.Header().Add(k, strings.Join(v, ", "))
		}
		w.Write([]byte("response write ok"))

	}))

	http.HandleFunc("/healthz", interceptor(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 ok"))
		// 4. 当访问 localhost/healthz 时，应返回 200
		w.WriteHeader(http.StatusOK)
	}))
	http.ListenAndServe(":3000", nil)
}
