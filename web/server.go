package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dafsic/go-pkgs/config"
	"github.com/dafsic/go-pkgs/mxlog"
	"github.com/dafsic/go-pkgs/web/middlewares"
	"github.com/dafsic/go-pkgs/web/middlewares/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Cfg struct {
	Addr    string // 192.168.1.100:8080
	AuthCfg auth.Cfg
}

type Server interface {
	RegisterHandler(method, path string, h gin.HandlerFunc)
}

type ServerImpl struct {
	listenAddr    string
	l             *mxlog.Logger
	srv           *http.Server
	gin           *gin.Engine
	authenticator auth.Authenticator
}

func NewServer(lc fx.Lifecycle, cfg config.Config, log mxlog.Loggers) Server {
	c := cfg.GetElem("server").(Cfg)
	impl := ServerImpl{
		l:          log.GetLogger("http"),
		listenAddr: c.Addr,
	}
	impl.authenticator = auth.NewAuthenticator(c.AuthCfg)
	impl.l.Info("Init...")

	impl.gin = gin.New()
	impl.gin.Use(middlewares.Record(impl.l))
	impl.gin.Use(middlewares.Cors())

	impl.srv = &http.Server{
		Addr:         impl.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      impl.gin,
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go func() {
				if err := impl.srv.ListenAndServe(); err != nil {
					impl.l.Error(err)
				}
			}()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			impl.srv.Shutdown(ctx)
			return nil
		},
	})

	return &impl
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
