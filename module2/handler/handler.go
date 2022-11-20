package handler

import (
	"net/http"
	"strings"
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
