package main

import (
	"log"
	"sync"
)

type LRUCache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type node struct {
	key   string
	value int
	pre   *node
	next  *node
}

type lruCache struct {
	cache    map[string]*node
	mutex    *sync.Mutex
	head     *node
	tail     *node
	capacity int
}

func (l *lruCache) Get(key string) (int, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currNode, ok := l.cache[key]
	if !ok {
		return 0, false
	}

	l.moveNodeToFront(currNode)

	return currNode.value, true
}

func (l *lruCache) moveNodeToFront(node *node) {
	if node == l.head {
		return
	}

	//Unlinked the node
	if node.next != nil {
		node.next.pre = node.pre
	}

	if node.pre != nil {
		node.pre.next = node.next
	}

	if l.head != nil {
		l.head.pre = node
		node.next = l.head
	}
	l.head = node

	if l.tail == nil {
		l.tail = node
	}
}

func (l *lruCache) evict() {
	if l.tail != nil {
		key := l.tail.key
		pre := l.tail.pre
		l.tail.pre = nil
		if pre != nil {
			pre.next = nil
		}
		l.tail = pre

		delete(l.cache, key)
	}
}

func (l *lruCache) Put(key string, val int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currNode, ok := l.cache[key]
	if ok {
		currNode.value = val
		l.moveNodeToFront(currNode)
		return
	}

	newNode := &node{
		key:   key,
		value: val,
	}

	l.cache[key] = newNode

	l.moveNodeToFront(newNode)

	if l.capacity < len(l.cache) {
		l.evict()
	}
}

func NewLRUCache(cap int) LRUCache {
	if cap <= 0 {
		panic("invalid cap")
	}

	return &lruCache{
		cache:    map[string]*node{},
		mutex:    &sync.Mutex{},
		capacity: cap,
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
