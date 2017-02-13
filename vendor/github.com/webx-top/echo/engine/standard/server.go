package standard

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Server struct {
		*http.Server
		config  *engine.Config
		handler engine.Handler
		logger  logger.Logger
		pool    *pool
	}

	pool struct {
		request        sync.Pool
		response       sync.Pool
		requestHeader  sync.Pool
		responseHeader sync.Pool
		url            sync.Pool
	}
)

func New(addr string) *Server {
	c := &engine.Config{Address: addr}
	return NewWithConfig(c)
}

func NewWithTLS(addr, certFile, keyFile string) *Server {
	c := &engine.Config{
		Address:     addr,
		TLSCertFile: certFile,
		TLSKeyFile:  keyFile,
	}
	return NewWithConfig(c)
}

func NewWithConfig(c *engine.Config) (s *Server) {
	s = &Server{
		Server: &http.Server{
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			Addr:         c.Address,
		},
		config: c,
		pool: &pool{
			request: sync.Pool{
				New: func() interface{} {
					return &Request{}
				},
			},
			response: sync.Pool{
				New: func() interface{} {
					return &Response{logger: s.logger}
				},
			},
			requestHeader: sync.Pool{
				New: func() interface{} {
					return &Header{}
				},
			},
			responseHeader: sync.Pool{
				New: func() interface{} {
					return &Header{}
				},
			},
			url: sync.Pool{
				New: func() interface{} {
					return &URL{}
				},
			},
		},
		handler: engine.ClearHandler(engine.HandlerFunc(func(req engine.Request, res engine.Response) {
			s.logger.Error("handler not set, use `SetHandler()` to set it.")
		})),
		logger: log.GetLogger("echo"),
	}
	s.Handler = s
	return
}

func (s *Server) SetHandler(h engine.Handler) {
	s.handler = engine.ClearHandler(h)
}

func (s *Server) SetLogger(l logger.Logger) {
	s.logger = l
}

// Start implements `engine.Server#Start` function.
func (s *Server) Start() error {
	if s.config.Listener == nil {
		return s.startDefaultListener()
	}
	return s.startCustomListener()
}

// Stop implements `engine.Server#Stop` function.
func (s *Server) Stop() error {
	if s.config.Listener == nil {
		return nil
	}
	return s.config.Listener.Close()
}

func (s *Server) startDefaultListener() error {
	ln, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}
	if s.config.TLSConfig != nil {
		s.logger.Info(`StandardHTTP is running at `, s.config.Address, ` [TLS]`)
		s.config.Listener = tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, s.config.TLSConfig)
	} else if len(s.config.TLSCertFile) > 0 && len(s.config.TLSKeyFile) > 0 {
		// TODO: https://github.com/golang/go/commit/d24f446a90ea94b87591bf16228d7d871fec3d92
		config := &tls.Config{}
		if !s.config.DisableHTTP2 {
			config.NextProtos = append(config.NextProtos, "h2")
		}
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
		if err != nil {
			return err
		}
		s.logger.Info(`StandardHTTP is running at `, s.config.Address, ` [TLS]`)
		s.config.Listener = tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	} else {
		s.logger.Info(`StandardHTTP is running at `, s.config.Address)
		s.config.Listener = tcpKeepAliveListener{ln.(*net.TCPListener)}
	}
	return s.Serve(s.config.Listener)
}

func (s *Server) startCustomListener() error {
	return s.Serve(s.config.Listener)
}

// ServeHTTP implements `http.Handler` interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Request
	req := s.pool.request.Get().(*Request)
	reqHdr := s.pool.requestHeader.Get().(*Header)
	reqHdr.reset(r.Header)
	reqURL := s.pool.url.Get().(*URL)
	reqURL.reset(r.URL)
	req.reset(r, reqHdr, reqURL)
	req.config = s.config

	// Response
	res := s.pool.response.Get().(*Response)
	resHdr := s.pool.responseHeader.Get().(*Header)
	resHdr.reset(w.Header())
	res.reset(w, r, resHdr)
	res.config = s.config

	s.handler.ServeHTTP(req, res)

	s.pool.request.Put(req)
	s.pool.requestHeader.Put(reqHdr)
	s.pool.url.Put(reqURL)
	s.pool.response.Put(res)
	s.pool.responseHeader.Put(resHdr)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
