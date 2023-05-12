package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dafsic/go-pkgs/config"
	"github.com/dafsic/go-pkgs/mxlog"
	"github.com/dafsic/go-pkgs/server/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Cfg struct {
	Addr string // 192.168.1.100:8080
}
type Server interface {
	RegisterHandler(method, path string, h gin.HandlerFunc)
}

type ServerImpl struct {
	listenAddr string
	l          *mxlog.Logger
	srv        *http.Server
	gin        *gin.Engine
}

func NewServer(lc fx.Lifecycle, c config.Config, log mxlog.Loggers) Server {
	cfg := c.GetElem("server").(Cfg)
	s := ServerImpl{
		l:          log.GetLogger("http"),
		listenAddr: cfg.Addr,
	}
	s.l.Info("Init...")

	s.gin = gin.New()
	s.gin.Use(middlewares.Record(s.l))

	s.srv = &http.Server{
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.gin,
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go func() {
				if err := s.srv.ListenAndServe(); err != nil {
					s.l.Error(err)
				}
			}()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			s.srv.Shutdown(ctx)
			return nil
		},
	})

	return &s
}

func (s *ServerImpl) RegisterHandler(method, path string, h gin.HandlerFunc) {
	switch strings.ToUpper(method) {
	case "GET":
		s.gin.GET(path, h)
	case "POST":
		s.gin.POST(path, h)
	case "PUT":
		s.gin.PUT(path, h)
	case "DELETE":
		s.gin.DELETE(path, h)
	default:
		//TODO: return error
	}
}

// Module for fx
var ServerModule = fx.Options(fx.Provide(NewServer))
