package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// 参考：https://gist.github.com/Boerworz/b683e46ae0761056a636
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
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

func GetResponseStatus(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)
		handler(lrw, r)
		glog.Infof("request %s status code %d", r.URL.Path, lrw.statusCode)

	}
}

func GetIP(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
		// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
		addr, err := getIP(r)
		if err != nil {
			glog.Errorf("get addr error %v\n", err)
		}
		glog.Infof("request %s ip addr %s\n", r.URL.Path, addr)

	}
}

func SetVersion(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		w.Header().Add("version", os.Getenv("VERSION"))
		handler(w, r)
	}
}
