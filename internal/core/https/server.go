package https

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/aob"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/metas"
	"github.com/mel0dys0ng/song/pkg/sys"
	"github.com/mel0dys0ng/song/pkg/tjme"
	"go.uber.org/zap"
)

type (
	Server struct {
		*Options

		engine     *gin.Engine
		httpServer *http.Server
		listener   net.Listener

		mt metas.MetadataInterface
	}
)

// New return new http server
// @Param opts []Option the option of http server
func New(opts []Option) *Server {
	return &Server{
		mt:      metas.Metadata(),
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
	// initServer gin
	s.engine = gin.New()
	isDebug := s.mt.Mode().IsModeDebug()
	gin.SetMode(aob.VarOrVar(isDebug, gin.DebugMode, gin.ReleaseMode))

	// use recover and trace middleware
	s.engine.Use(s.setupRecoverAndTraceMiddleware())

	// use client middleware
	s.engine.Use(s.setupClientMiddleware())

	//use cors middleware
	s.engine.Use(s.setupCORSMiddleware())

	// use csrf middleware
	s.engine.Use(s.setupCSRFMiddleware())

	// use sign middleware
	s.engine.Use(s.setupSignMiddleware())

	// use responded middleware
	s.engine.Use(s.setupRespondedMiddleware())

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
				sys.Panic(err.Error())
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
		if err := s.listen(); err != nil {
			erlogs.Convert(err).Options(BaseELOptions()).RecordLog(ctx)
			return
		}

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

	erlogs.New("server exiting...").Options(BaseELOptions()).InfoLog(ctx)
}

func (s *Server) listen() (err error) {
	fields := erlogs.OptionFields(zap.String("addr", s.Addr))

	// 设置监听器的监听对象（新建的或已存在的 socket 描述符）
	s.listener, err = net.Listen("tcp", s.Addr)

	// 监听失败
	if err != nil {
		erlogs.Convert(err).Wrap("listen server failed").Options(BaseELOptions()).PanicLog(context.Background(), fields)
		return
	}

	// 监听成功
	erlogs.New("listen server success").Options(BaseELOptions()).InfoLog(context.Background(), fields)

	return
}

func (s *Server) serve() {
	if s.OnStart != nil {
		s.OnStart()
	}

	var err error
	ctx := context.Background()
	fields := erlogs.OptionFields(zap.Int("pid", os.Getpid()), zap.String("addr", s.Addr))

	// serve
	if s.TLSOpen {
		err = s.httpServer.ServeTLS(s.listener, s.TLSCertFile, s.TLSKeyFile)
	} else {
		err = s.httpServer.Serve(s.listener)
	}

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
		erlogs.New("server closed").Options(BaseELOptions()).InfoLog(ctx, fields)
	}

	// serve failed
	if err != nil {
		if s.OnStartFail != nil {
			s.OnStartFail(err)
		}
		erlogs.Convert(err).Wrap("failed to serve").Options(BaseELOptions()).PanicLog(ctx, fields)
		return
	}

	// serve success
	erlogs.New("serve server success").Options(BaseELOptions()).InfoLog(ctx, fields)
}

func (s *Server) shutdown() {
	hammerTime := tjme.ParseDuration(s.HammerTime, DefaultHammerTime)
	ctx, cancel := context.WithTimeout(context.Background(), hammerTime)
	defer cancel()

	fields := erlogs.OptionFields(zap.Int("pid", os.Getpid()), zap.String("addr", s.Addr))
	if err := s.httpServer.Shutdown(ctx); err != nil {
		erlogs.Convert(err).Wrap("failed to shutdown").Options(BaseELOptions()).PanicLog(ctx, fields)
		return
	}

	erlogs.New("shutdown server success").Options(BaseELOptions()).InfoLog(ctx, fields)
}
