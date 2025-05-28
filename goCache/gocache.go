package gocache

import (
	"fmt"
	"log"
	"sync"
)

// a getter loads data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// a getter func implements getter with a function
type Getterfunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f Getterfunc) Get(key string) ([]byte, error) {
	return f(key)
}

// a group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// global variables
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// newgroup create a new instance of group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there is no such group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// get value for a key from cache
func (g *Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[GoCache] hit")
		return v, nil
	}
	return g.load(key)
}

// if key not exist in cache, then load would use getLocally to callback data
func (g *Group) load(key string) (value Byteview, err error) {
	return g.getLocally(key)
}

// this is the callback function to retrieve data not in lru cache
func (g *Group) getLocally(key string) (Byteview, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return Byteview{}, err
	}

	value := Byteview{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// the data from callback will be added back to the lru cache
func (g *Group) populateCache(key string, value Byteview) {
	g.mainCache.add(key, value)
}
