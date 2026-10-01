package main

import (
	"container/list"
	"log"
	"sync"
)

type node2 struct {
	key string
	val int
}

type lruCache struct {
	cache    map[string]*list.Element
	capacity int
	list     *list.List
	mutex    *sync.Mutex
}

func (l *lruCache) Get(key string) (int, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currNode, ok := l.cache[key]
	if !ok {
		return 0, false
	}

	//Push item to front.
	l.list.MoveToFront(currNode)

	val := currNode.Value.(*node2)

	return val.val, true
}

func (l *lruCache) Put(key string, val int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currEle, ok := l.cache[key]
	if ok {
		currEle.Value = &node2{key: key, val: val}

		l.list.MoveToFront(currEle)
	}

	newNode := &node2{key: key, val: val}
	ele := l.list.PushFront(newNode)
	l.cache[key] = ele

	if len(l.cache) > l.capacity {
		l.list.Remove(l.list.Back())
	}
}

type LRUCache interface {
	Get(string) (int, bool)
	Put(string, int)
}

func NewLRUCache(cap int) LRUCache {
	if cap <= 0 {
		panic("invalid cap")
	}

	return &lruCache{
		cache:    map[string]*list.Element{},
		capacity: cap,
		list:     list.New(),
		mutex:    &sync.Mutex{},
	}
}

func main() {
	lru := NewLRUCache(3)

	lru.Put("1", 100)
	lru.Put("2", 200)
	lru.Put("3", 300)

	log.Println(lru.Get("1"))
	log.Println(lru.Get("2"))
	log.Println(lru.Get("3"))
	log.Println(lru.Get("4"))

	lru.Put("4", 400)

	log.Println(lru.Get("4"))

	log.Println(lru.Get("1"))
}
