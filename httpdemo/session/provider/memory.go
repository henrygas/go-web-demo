package memory

import (
	"container/list"
	"errors"
	"go-web/httpdemo/session"
	"sync"
	"time"
)

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	values       map[interface{}]interface{}
}

func (ss *SessionStore) Set(key, value interface{}) error {
	ss.values[key] = value
	_ = provider.SessionUpdate(ss.sid)
	return nil
}

func (ss *SessionStore) Get(key interface{}) interface{} {
	_ = provider.SessionUpdate(ss.sid)
	if value, ok := ss.values[key]; ok {
		return value
	} else {
		return nil
	}
}

func (ss *SessionStore) Delete(key interface{}) error {
	if _, ok := ss.values[key]; ok {
		delete(ss.values, key)
	}
	_ = provider.SessionUpdate(ss.sid)
	return nil
}

func (ss *SessionStore) SessionID() (string, error) {
	return ss.sid, nil
}

type Provider struct {
	lock     sync.Mutex               // lock
	sessions map[string]*list.Element // sid -> session element
	list     *list.List               // session element queue, used for GC
}

var provider = &Provider{list: list.New()}

// create a new session, and push it to front of session element list
func (pd *Provider) SessionInit(sid string) (session.Session, error) {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsession := &SessionStore{
		sid:          sid,
		timeAccessed: time.Now(),
		values:       v,
	}
	element := pd.list.PushFront(newsession)
	pd.sessions[sid] = element
	return newsession, nil
}

// read session from provider if session exists,
// or else return nil
func (pd *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pd.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		return nil, errors.New("session not exits")
	}
}

// update session, and move it to the front of session element list
func (pd *Provider) SessionUpdate(sid string) error {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	if element, ok := pd.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pd.list.MoveToFront(element)
	}
	return nil
}

// delete session from provider, from both map and list
func (pd *Provider) SessionDestroy(sid string) error {
	if element, ok := pd.sessions[sid]; ok {
		delete(pd.sessions, sid)
		pd.list.Remove(element)
	}
	return nil
}

// remove expired session elements from both map and list
func (pd *Provider) SessionGC(maxLifeTime int64) {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	for {
		element := pd.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			delete(pd.sessions, element.Value.(*SessionStore).sid)
			pd.list.Remove(element)
		} else {
			break
		}
	}
}

func init() {
	provider.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", provider)
}
