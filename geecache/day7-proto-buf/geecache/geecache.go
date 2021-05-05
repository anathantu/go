package geecache

import (
	"example1/geecache/geecachepb"
	"example1/geecache/singleflight"
	"fmt"
	"log"
	"sync"
)

/**
*一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。
*比如可以创建三个 Group，缓存学生的成绩命名为 scores，缓存学生信息的命名为 info，
*缓存学生课程的命名为 courses。
 */

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker

	//使用singleflight来保证多个请求只有一个请求能够落实运行
	loader *singleflight.Group
}

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 定义了一个函数类型
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

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
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(key string) *Group {
	mu.RLock()
	g := groups[key]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	//没有命中缓存时加载缓存
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
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
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &geecachepb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &geecachepb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	//调用Group中的回调函数
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
