package main

import (
	"log"
	"sync"
)

//LRU is a Cache eviction Policy where the key that is used last is evicted

// At each Get or PUT the key is put at the front.
// Only value farthest accessed is removed.

type LRUCache interface {
	Get(string) int
	Put(string, int)
}

type node struct {
	val  int
	key  string
	pre  *node
	next *node
}

type lruCache struct {
	cache    map[string]*node
	head     *node
	tail     *node
	capacity int
	mutex    *sync.Mutex
}

func NewLRUCache(cap int) LRUCache {
	if cap <= 0 {
		panic("invalid capacity")
	}

	return &lruCache{
		cache:    make(map[string]*node),
		head:     nil,
		tail:     nil,
		capacity: cap,
		mutex:    &sync.Mutex{},
	}
}

func (l *lruCache) Get(key string) int {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	node, ok := l.cache[key]
	if !ok {
		return -1
	}

	l.moveToFront(node)

	return node.val

}

func (l *lruCache) moveToFront(node *node) {
	if node == l.head {
		return
	}

	if node.next != nil {
		node.next.pre = node.pre
	}
	if node.pre != nil {
		node.pre.next = node.next
	}

	if node == l.tail {
		l.tail = node.pre
	}

	node.pre = nil
	if l.head != nil {
		l.head.pre = node
		node.next = l.head
	}
	l.head = node

	if l.tail == nil {
		l.tail = node
	}
}

func (l *lruCache) evictTail() {
	if l.tail != nil {
		key := l.tail.key
		if l.tail.pre != nil {
			l.tail.pre.next = nil
		}
		l.tail = l.tail.pre
		delete(l.cache, key)
	}
}

func (l *lruCache) Put(key string, value int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currVal, ok := l.cache[key]
	if ok {
		currVal.val = value
		l.moveToFront(currVal)
		return
	}

	newNode := &node{
		val:  value,
		pre:  nil,
		next: nil,
		key:  key,
	}
	l.moveToFront(newNode)
	l.cache[key] = newNode

	if len(l.cache) > l.capacity {
		l.evictTail()
	}
}

func main() {
	lru := NewLRUCache(5)
	lru.Put("1", 100)
	lru.Put("2", 200)
	lru.Put("3", 300)
	lru.Put("4", 400)
	lru.Put("5", 500)

	log.Println(lru.Get("1"))
	log.Println(lru.Get("2"))
	log.Println(lru.Get("3"))
	log.Println(lru.Get("4"))
	log.Println(lru.Get("5"))

	log.Println(lru.Get("6"))

	lru.Put("6", 600)

	log.Println(lru.Get("1"))
}
