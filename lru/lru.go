package lru

import (
	"container/list"
	"sync"
)

type Cache interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	Remove(key string)
	Refresh() Cache
}

//原生双向链表实现LRU
type LRUCache struct {
	m        map[string]*Node
	l        *list.List
	capacity int
	mu      sync.Mutex
}

type Node struct {
	Path string
	e *list.Element
}

func Constructor(capacity int) Cache {
	return &LRUCache{
		m:        map[string]*Node{},
		l:        list.New(),
		capacity: capacity,
	}
}

func (this *LRUCache) Get(key string) interface{} {
	if value, ok := this.m[key]; ok {
		this.l.MoveToFront(value.e)
		return value.Path
	} else {
		return ""
	}
}

func (this *LRUCache) Put(key string, value interface{}) {
	path := value.(string)
	if nodeV, ok := this.m[key]; ok {
		this.l.MoveToFront(nodeV.e)
		this.m[key].Path = path
	} else {
		this.l.PushFront(key)
		this.m[key] = &Node{
			Path: path,
			e: this.l.Front(),
		}
	}
	if len(this.m) > this.capacity { //容量到达最大
		delete(this.m, this.l.Remove(this.l.Back()).(string))
	}
}

func (this *LRUCache) Remove(key string) {
	nodeV, ok := this.m[key]
	if !ok {
		return
	}
	this.l.Remove(nodeV.e)
	delete(this.m, key)
}

func (this *LRUCache) Refresh() Cache{
	return &LRUCache{
		m:        map[string]*Node{},
		l:        list.New(),
		capacity: this.capacity,
	}
}