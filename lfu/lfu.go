package lfu

import (
	"container/heap"

	"github.com/iscreen/go-cache"
)

type lfu struct {
	maxBytes  int
	onEvicted func(key string, value interface{})
	usedBytes int
	queue     *queue
	cache     map[string]*entry
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

// Get 從 cache 中獲取 key 對應的值, nil 表示不存在
func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}

	return nil
}

// Del 從 cache 中刪除 key 對應的元素
func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

// DelOldest 以 cache 中刪除最舊的記錄
func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

// Len 返回當前的 cache 記錄數
func (l *lfu) Len() int {
	return l.queue.Len()
}

func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
	}
	en := &entry{key: key, value: value}
	heap.Push(l.queue, en)
	l.cache[key] = en

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeElement(heap.Pop(l.queue))
	}
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}

	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}
