package session

import (
	"sync"
	"errors"
	"net/http"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"time"
)

type Session interface {
	Id() string
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Delete(key interface{})
}

type Provider interface {
	SessionStart(sid string) Session
	SessionEnd(sid string)
	Gc(maxLifeTime int64)
}

const (
	SESSION_PROVIDER_MEMORY = "memory"
)

var providers = make(map[string]Provider, 1)

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	maxLifeTime int64
	provider    Provider
}

func RegisterProvider(name string, provider Provider) {
	if provider == nil {
		panic("provider must not be nil")
	}
	if _, ok := providers[name]; ok {
		panic("dupulicated register provider " + name)
	}
	providers[name] = provider
}

func NewSessionManager(cookieName, providerName string, maxLifeTime int64) (*SessionManager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, errors.New("unknown provider is named " + providerName)
	}
	mgr := &SessionManager{
		cookieName:cookieName,
		maxLifeTime:maxLifeTime,
		provider:provider,
	}
	go mgr.gc()
	return mgr, nil
}

func newSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (mgr *SessionManager)SessionStart(w http.ResponseWriter, r *http.Request) Session {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	cookie, err := r.Cookie(mgr.cookieName)
	if err != nil || cookie.Value == "" {
		sid := newSessionId()
		if sid != "" {
			session := mgr.provider.SessionStart(sid)
			if session != nil {
				http.SetCookie(w, &http.Cookie{
					Name:mgr.cookieName,
					Value:sid,
					Path:"/",
					HttpOnly:true,
					MaxAge:int(mgr.maxLifeTime),
				})
			}
			return session
		}

	} else {
		sid := url.PathEscape(cookie.Value)
		session := mgr.provider.SessionStart(sid)
		if session != nil {
			return session
		}
	}
	return nil
}

func (mgr *SessionManager)SessionEnd(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(mgr.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.provider.SessionEnd(cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:mgr.cookieName,
		Path:"/",
		HttpOnly:true,
		Expires:time.Now(),
		MaxAge:-1,
	})
}

func (mgr *SessionManager)gc() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.provider.Gc(mgr.maxLifeTime)
	time.AfterFunc(time.Duration(mgr.maxLifeTime), func() {
		mgr.gc()
	})
}
