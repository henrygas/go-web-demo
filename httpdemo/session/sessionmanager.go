package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() (string, error)
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxlifetime int64)
	SessionUpdate(sid string) error
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
}

func (sm *Manager) SessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (sm *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil || cookie.Value == "" {
		sid := sm.SessionID()
		session, _ = sm.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:     sm.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(sm.maxLifeTime),
		}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = sm.provider.SessionRead(sid)
		if session == nil {
			session, _ = sm.provider.SessionInit(sid)
			cookie := http.Cookie{
				Name:     sm.cookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(sm.maxLifeTime),
			}
			http.SetCookie(w, &cookie)
		}
	}
	return session
}

func (sm *Manager) GetSession(r *http.Request) (Session, error) {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil || cookie.Value == "" {
		return nil, err
	} else {
		return sm.provider.SessionRead(cookie.Value)
	}
}

func (sm *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		sm.lock.Lock()
		defer sm.lock.Unlock()
		_ = sm.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{
			Name:     sm.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1,
		}
		http.SetCookie(w, &cookie)
	}
}

func (sm *Manager) GC() {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.provider.SessionGC(sm.maxLifeTime)
	time.AfterFunc(time.Duration(sm.maxLifeTime), func() {
		sm.GC()
	})
}

var provides = make(map[string]Provider)

func NewManager(providerName string, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q (forgotten import?)", providerName)
	}
	return &Manager{
		cookieName:  cookieName,
		provider:    provider,
		maxLifeTime: maxLifeTime,
	}, nil
}

func Register(providerName string, provider Provider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[providerName]; dup {
		panic("session: Register called twice for provider " + providerName)
	}
	provides[providerName] = provider
}
