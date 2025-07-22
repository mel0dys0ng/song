package internal

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/song/erlogs"
	"github.com/song/metas"
	"github.com/song/utils/aob"
	"github.com/song/utils/systems"
	"github.com/song/utils/tjme"
	"go.uber.org/zap"
)

var (
	_elgSys     erlogs.ErLogInterface
	_elgSysOnce sync.Once
)

type (
	Server struct {
		*Options

		engine     *gin.Engine
		httpServer *http.Server
		listener   net.Listener
	}
)

// New return new http server
// @Param opts []Option the option of http server
func New(opts []Option) *Server {
	return &Server{
		Options: newOptions(opts),
	}
}

// Serve 启动服务
func (s *Server) Serve() {
	s.loadDefers()
	s.initServer()
	s.loadInits()
	s.loadMiddlewares()
	s.loadRoutes()
	s.runServer()
}

func (s *Server) initServer() {
	// initServer erlog
	_elgSysOnce.Do(func() {
		_elgSys = erlogs.New(
			erlogs.Log(true),
			erlogs.TypeSystem(),
			erlogs.Logger(s.Config.ErLog),
			erlogs.Msgf("[https] %s"),
		)
	})

	// initServer gin
	s.engine = gin.New()
	isDebug := metas.Mode().IsModeTestOrPreOrDebug()
	gin.SetMode(aob.Aorb(isDebug, gin.DebugMode, gin.ReleaseMode))

	s.engine.Use(gin.Recovery())
	s.engine.Use(s.trace)
	s.engine.Use(s.responded)

	if s.Cors != nil && s.Cors.Enable {
		s.engine.Use(cors.New(cors.Config{
			AllowOrigins:     s.Cors.AllowOrigins,
			AllowMethods:     s.Cors.AllowMethods,
			AllowHeaders:     s.Cors.AllowHeaders,
			AllowCredentials: s.Cors.AllowCredentials,
			AllowWildcard:    s.Cors.AllowWildcard,
			ExposeHeaders:    s.Cors.ExposeHeaders,
			MaxAge:           tjme.ParseDuration(s.Cors.MaxAge, DefaultCorsMaxAge),
		}))
	}

	// initServer http server
	s.httpServer = &http.Server{
		Handler:           s.engine,
		Addr:              s.Addr,
		ReadTimeout:       tjme.ParseDuration(s.ReadTimeout, DefaultReadTimeout),
		ReadHeaderTimeout: tjme.ParseDuration(s.ReadHeaderTimeout, DefaultReadHeaderTimeout),
		WriteTimeout:      tjme.ParseDuration(s.WriteTimeout, DefaultWriteTimeout),
		IdleTimeout:       tjme.ParseDuration(s.IdleTimeout, DefaultIdleTimeout),
		MaxHeaderBytes:    s.MaxHeaderBytes,
	}

	s.httpServer.SetKeepAlivesEnabled(s.KeepAlive)
	if s.OnShutdown != nil {
		s.httpServer.RegisterOnShutdown(s.OnShutdown)
	}
}

func (s *Server) loadInits() {
	if len(s.Inits) > 0 {
		defer func() { s.Inits = nil }()
		for _, fn := range s.Inits {
			if err := fn(); err != nil {
				systems.Panic(err.Error())
			}
		}
	}
}

func (s *Server) loadDefers() {
	if len(s.Defers) > 0 {
		defer func() {
			for _, fn := range s.Defers {
				fn()
			}
		}()
	}
}

func (s *Server) loadRoutes() {
	if len(s.Routes) > 0 {
		defer func() { s.Routes = nil }()
		for _, fn := range s.Routes {
			fn(s.engine)
		}
	}
}

func (s *Server) loadMiddlewares() {
	if len(s.Middlewares) > 0 {
		defer func() { s.Middlewares = nil }()
		sort.Slice(s.Middlewares, func(i, j int) bool {
			return s.Middlewares[i].Priority <= s.Middlewares[j].Priority
		})
		for _, v := range s.Middlewares {
			fc := v.Handle(s.engine)
			s.engine.Use(fc)
		}
	}
}

func (s *Server) runServer() {
	// 创建一个带有信号的 context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 服务
	go func() {
		// 监听服务
		s.listen()
		// 启动服务
		s.serve()
	}()

	// 等待 context 被取消（即收到信号）
	<-ctx.Done()

	// 关闭
	s.shutdown()

	// 执行退出回调
	if s.OnExit != nil {
		s.OnExit()
	}

	// 退出
	_elgSys.InfoL(context.Background(),
		erlogs.Msgv("server exiting..."),
		erlogs.Fields(zap.String("addr", s.Addr)),
	)

	os.Exit(0)
}

func (s *Server) listen() {
	fields := erlogs.Fields(zap.String("addr", s.Addr))

	var err error
	// 设置监听器的监听对象（新建的或已存在的 socket 描述符）
	s.listener, err = net.Listen("tcp", s.Addr)
	// 监听失败
	if err != nil {
		_elgSys.PanicL(
			context.Background(),
			erlogs.Msgv("listen server failed"),
			erlogs.Content(err.Error()),
			fields,
		)
		return
	}

	// 监听成功
	_elgSys.InfoL(context.Background(), erlogs.Msgv("listen server success"), fields)
}

func (s *Server) serve() {
	if s.OnStart != nil {
		s.OnStart()
	}

	fields := erlogs.Fields(zap.Int("pid", os.Getpid()), zap.String("addr", s.Addr))

	// serve
	var err error
	if s.TLSOpen {
		err = s.httpServer.ServeTLS(s.listener, s.TLSCertFile, s.TLSKeyFile)
	} else {
		err = s.httpServer.Serve(s.listener)
	}

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
		_elgSys.InfoL(context.Background(), erlogs.Msgv("server closed"), fields)
	}

	// serve failed
	if err != nil {
		if s.OnStartFail != nil {
			s.OnStartFail(err)
		}
		_elgSys.PanicL(
			context.Background(),
			erlogs.Msgv("failed to serve server"),
			erlogs.Content(err.Error()),
			fields,
		)
		return
	}
}

func (s *Server) shutdown() {
	hammerTime := tjme.ParseDuration(s.HammerTime, DefaultHammerTime)
	ctx, cancel := context.WithTimeout(context.Background(), hammerTime)
	defer cancel()

	fields := erlogs.Fields(zap.Int("pid", os.Getpid()), zap.String("addr", s.Addr))
	if err := s.httpServer.Shutdown(ctx); err != nil {
		_elgSys.PanicL(
			context.Background(),
			erlogs.Msgv("failed to shutdown server"),
			erlogs.Content(err.Error()),
			fields,
		)
		return
	}

	_elgSys.InfoL(context.Background(), erlogs.Msgv("shutdown server success"), fields)
}
