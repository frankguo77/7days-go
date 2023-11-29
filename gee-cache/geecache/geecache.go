package geecache

import (
	"fmt"
	"geecache/singleflight"
	 pb "geecache/geecachepb"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache  cache
	peers      PeerPicker
	loader    *singleflight.Group
}

type Getter interface {
	Get(key string) ([]byte ,error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte ,error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeers called more than once.")
	}

	g.peers = peers
}

func (g* Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value Byteview, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}

				log.Println("[GeeCache] Failed to get from peer", err)
			}

		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(Byteview), nil
	}

	return 
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (Byteview, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}

	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return Byteview{}, err
	}

	return Byteview{b: res.Value}, nil
}

func (g *Group) getLocally(key string) (Byteview, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return Byteview{}, err
	}

	value := Byteview{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value Byteview) {
	g.mainCache.add(key, value)
}
