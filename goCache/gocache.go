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
	peers     PeerPicker
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
	// 使用 PickPeer() 方法选择节点，若非本机节点，则调用 getFromPeer() 从远程获取。若是本机节点或失败，则回退到 getLocally()。
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GoCache] Failed to get from peer", err)
		}
	}

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

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	// 新增 RegisterPeers() 方法，将 实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中。
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (Byteview, error) {
	// 新增 getFromPeer() 方法，使用实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值。
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return Byteview{}, err
	}
	return Byteview{b: bytes}, nil
}
