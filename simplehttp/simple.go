package simplehttp

import (
	"net/http"
	"github.com/zuowzhang/go-support/session"
	"net"
	"time"
	"sync"
)

type (
	FilterFunc func(http.HandlerFunc) http.HandlerFunc

	HandlerFunc func(Context) error

	Simple struct {
		sessionManager *session.SessionManager;
		server         *http.Server
		listener       net.Listener
		filters        []FilterFunc
		router         *Router
		pool           sync.Pool
	}
)

func NewSimple() *Simple {
	simple := new(Simple)
	manager, err := session.NewSessionManager("simple-cookie", session.SESSION_PROVIDER_MEMORY, 3600);
	if err == nil {
		simple.sessionManager = manager
	}
	simple.server = new(http.Server)
	simple.server.Handler = simple
	simple.pool.New = func() interface{} {
		return simple.newContext()
	}
	return simple
}

func (s *Simple)newContext() Context {
	return &context{
		simple:s,
	}
}

func (s *Simple)Use(filter ...FilterFunc) *Simple {
	s.filters = append(s.filters, filter...)
	return s
}

func (s *Simple)Get(path string, h HandlerFunc) {
	s.router.Add(http.MethodGet, path, h)
}

func (s *Simple)Post(path string, h HandlerFunc) {
	s.router.Add(http.MethodPost, h)
}

func (s *Simple)Start(address string) error {
	s.server.Addr = address
	return s.StartServer(s.server)
}

func (s *Simple)StartServer(server *http.Server) (err error) {
	server.Handler = s
	if s.listener == nil {
		s.listener, err = newListener(server.Addr)
		if err != nil {
			return
		}
	}
	err = server.Serve(s.listener)
	return
}

func (s *Simple)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := s.pool.Get().(Context)
	defer s.pool.Put(context)
	context.Reset(w, r)

}

type tcpKeepAliveListener struct {
	listener *net.TCPListener
}

func (l tcpKeepAliveListener)Accept() (net.Conn, error) {
	tc, err := l.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute);
	return tc, err
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{
		listener:l,
	}, nil
}


