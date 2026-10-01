package main

import (
	"container/list"
	"log"
	"sync"
)

type LFUCache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type lfuCache struct {
	cache   map[string]*list.Element
	freq    map[int]*list.List
	cap     int
	mutex   *sync.Mutex
	minFreq int
}

type node struct {
	key  string
	val  int
	freq int
}

func (l *lfuCache) Get(key string) (int, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currEle, ok := l.cache[key]
	if !ok {
		return 0, false
	}
	val := currEle.Value.(*node)
	dl := l.freq[val.freq]
	dl.Remove(currEle)

	if val.freq == l.minFreq && dl.Len() == 0 {
		l.minFreq += 1
	}

	val.freq += 1

	if currList, ok := l.freq[val.freq]; ok {
		ele := currList.PushFront(val)
		l.cache[key] = ele
	} else {
		tempList := list.New()
		ele := tempList.PushFront(val)
		l.cache[key] = ele
		l.freq[val.freq] = tempList
	}

	return val.val, true
}

func (l *lfuCache) Put(key string, value int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	currEle, ok := l.cache[key]
	if ok {
		val := currEle.Value.(*node)
		dl := l.freq[val.freq]
		dl.Remove(currEle)

		if val.freq == l.minFreq && dl.Len() == 0 {
			l.minFreq += 1
		}

		val.freq += 1
		val.val = value

		if currList, ok := l.freq[val.freq]; ok {
			ele := currList.PushFront(val)
			l.cache[key] = ele
		} else {
			tempList := list.New()
			ele := tempList.PushFront(val)
			l.cache[key] = ele
			l.freq[val.freq] = tempList
		}

		return
	}

	if len(l.cache) >= l.cap {
		minDL := l.freq[l.minFreq]
		back := minDL.Back()
		v := back.Value.(*node)
		minDL.Remove(minDL.Back())
		delete(l.cache, v.key)
	}

	newNode := &node{
		key:  key,
		val:  value,
		freq: 1,
	}
	if currList, ok := l.freq[1]; ok {
		ele := currList.PushFront(newNode)
		l.cache[key] = ele
	} else {
		tempList := list.New()
		ele := tempList.PushFront(newNode)
		l.freq[1] = tempList
		l.cache[key] = ele
	}

	l.minFreq = 1
}

func NewLFU(cap int) LFUCache {
	return &lfuCache{
		cache:   make(map[string]*list.Element),
		freq:    make(map[int]*list.List),
		cap:     cap,
		mutex:   &sync.Mutex{},
		minFreq: 0,
	}
}

func main() {
	lru := NewLFU(5)
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
	log.Println(lru.Get("2"))
	log.Println(lru.Get("3"))
	log.Println(lru.Get("4"))
	log.Println(lru.Get("5"))
}
