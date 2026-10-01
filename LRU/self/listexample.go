package main

import "container/list"

func main() {

	l := list.New()
	e1 := l.PushFront(&node{1, 10}) // [1]
	e2 := l.PushFront(&node{2, 20}) // [2,1]
	l.MoveToFront(e1)               // [1,2]
	n := l.Back().Value.(*node)     // oldest -> node 2
	l.Remove(l.Back())              // [1]
	for e := l.Front(); e != nil; e = e.Next() {
		_ = e.Value.(*node)
	}

}

type node struct{ key, val int }
