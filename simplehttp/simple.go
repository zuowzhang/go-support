package simplehttp

import (
	"net/http"
	"github.com/zuowzhang/go-support/session"
	"net"
	"time"
	"sync"
)

type (
	HandlerFunc func(Context) error

	FilterFunc func(HandlerFunc) HandlerFunc

	Simple struct {
		sessionManager *session.SessionManager;
		server         *http.Server
		listener       net.Listener
		filters        []FilterFunc
		pool           sync.Pool
		router         *Router
		render         Render
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
	simple.router = NewRouter()
	simple.render = new(render)
	return simple
}

func (s *Simple)newContext() Context {
	return &context{
		simple:s,
	}
}

func (s *Simple)Group(prefix string, filters... FilterFunc) *Group {
	return &Group{
		prefix:prefix,
		filters:filters,
		simple:s,
	}
}

func (s *Simple)Use(filter ...FilterFunc) *Simple {
	s.filters = append(s.filters, filter...)
	return s
}

func (s *Simple)Get(path string, h HandlerFunc, filters... FilterFunc) {
	s.router.Add(http.MethodGet, path, h, filters...)
}

func (s *Simple)Post(path string, h HandlerFunc, filters... FilterFunc) {
	s.router.Add(http.MethodPost, path, h, filters...)
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

func (s *Simple)ParseGlob(pattern string) error {
	if s.render == nil {
		s.render = new(render)
	}
	return s.render.(*render).ParseGlob(pattern)
}

func (s *Simple)HttpErrorHandler(err error, c Context) {
	c.String(500, err.Error())
}

func (s *Simple)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := s.pool.Get().(Context)
	defer s.pool.Put(context)
	context.Reset(w, r)
	h := context.Handler();
	for i := len(s.filters) - 1; i >= 0; i-- {
		h = s.filters[i](h)
	}
	if err := h(context); err != nil {
		s.HttpErrorHandler(err, context)
	}
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (l tcpKeepAliveListener)Accept() (net.Conn, error) {
	tc, err := l.AcceptTCP()
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
		l.(*net.TCPListener),
	}, nil
}


