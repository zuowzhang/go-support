package session

import (
	"time"
	"sync"
	"container/list"
)

var provider = &memoryProvider{
	sessions:make(map[string]*list.Element, 0),
	sessionList:list.New(),
}

func init() {
	RegisterProvider(SESSION_PROVIDER_MEMORY, provider)
}

type memoryStore struct {
	id               string
	lastAccessedTime time.Time
	values           map[interface{}]interface{}
}

func (store *memoryStore)Id() string {
	return store.id
}

func (store *memoryStore)Get(key interface{}) interface{} {
	if value, ok := store.values[key]; ok {
		return value
	}
	return nil
}

func (store *memoryStore)Set(key, value interface{}) {
	store.values[key] = value
	provider.update(store)
}

func (store *memoryStore)Delete(key interface{}) {
	delete(store.values, key)
	provider.update(store)
}

type memoryProvider struct {
	lock        sync.RWMutex
	sessions    map[string]*list.Element
	sessionList *list.List
}

func (p *memoryProvider)update(store *memoryStore) {
	store.lastAccessedTime = time.Now()
	p.sessionList.MoveToFront(p.sessions[store.id])
}

func (p *memoryProvider)SessionStart(sid string) Session {
	p.lock.RLock()
	element, ok := p.sessions[sid]
	p.lock.RUnlock()
	if ok {
		return element.Value.(*memoryStore)
	}
	//new session
	p.lock.Lock()
	defer p.lock.Unlock()
	element = p.sessionList.PushFront(&memoryStore{
		id:sid,
		lastAccessedTime:time.Now(),
		values:make(map[interface{}]interface{}),
	})
	p.sessions[sid] = element
	return element.Value.(*memoryStore)
}

func (p *memoryProvider)SessionEnd(sid string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	element, ok := p.sessions[sid]
	if ok {
		delete(p.sessions, sid)
		p.sessionList.Remove(element)
	}
}

func (p *memoryProvider)Gc(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for {
		element := p.sessionList.Back()
		if element == nil {
			break
		}
		store := element.Value.(*memoryStore)
		if store.lastAccessedTime.Unix() + maxLifeTime < time.Now().Unix() {
			delete(p.sessions, store.id)
			p.sessionList.Remove(element)
		} else {
			break
		}
	}
}
