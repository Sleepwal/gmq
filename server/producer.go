package server

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type Producer interface {
	Start()
}

type HttpProducer struct {
	listenAddr  string
	server      *Server
	produceChan chan Message
}

func NewHttpProducer(listenAddr string, produceChan chan Message) *HttpProducer {
	return &HttpProducer{
		listenAddr:  listenAddr,
		produceChan: produceChan,
	}
}

func (p *HttpProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path  = strings.TrimPrefix(r.URL.Path, "/") //去掉URL路径的前缀
		parts = strings.Split(path, "/")            //将路径分割成parts数组
	)

	if r.Method == "POST" {
		if len(parts) != 3 {
			slog.Warn("invalid request")
			return
		}

		p.produceChan <- Message{
			Data:  []byte(parts[2]),
			Topic: parts[1],
		}
	}

	log.Println("parts: ", parts)
}

func (p *HttpProducer) Start() {
	slog.Info("starting producer HTTP server, address: ", p.listenAddr)
	err := http.ListenAndServe(p.listenAddr, p)
	if err != nil {
		slog.Warn("failed to start producer HTTP server", err)
	}
}
