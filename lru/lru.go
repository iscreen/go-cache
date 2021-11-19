package lru

import (
	"container/list"

	"github.com/iscreen/go-cache"
)

type lru struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	ll *list.List

	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := l.ll.PushBack(en)
	l.cache[key] = e
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.ll.Remove(e)

	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

// // New 創建一個新的 cache，如果 maxBytes 是 0，表示沒有容量限制
// func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
// 	return &lru{
// 		maxBytes:  maxBytes,
// 		onEvicted: onEvicted,
// 		ll:        list.New(),
// 		cache:     make(map[string]*list.Element),
// 	}
// }

// // Set 往 cache 尾部增加一個元素 如果 key 已存在，則放件尾部，並更新值
// func (l *lru) Set(key string, value interface{}) {
// 	if e, ok := l.cache[key]; ok {
// 		l.ll.MoveToBack(e)
// 		en := e.Value.(*entry)
// 		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
// 		en.value = value
// 		return
// 	}

// 	en := &entry{key, value}
// 	e := l.ll.PushBack(en)
// 	l.cache[key] = e

// 	l.usedBytes += en.Len()
// 	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
// 		l.DelOldest()
// 	}
// }

// // Get 從 cache 中獲取 key 對應的值，nil 表示 key 不存在
// func (l *lru) Get(key string) interface{} {
// 	if e, ok := l.cache[key]; ok {
// 		l.ll.MoveToBack(e)
// 		return e.Value.(*entry).value
// 	}

// 	return nil
// }

// // Del 從 cache 中刪除 key 對應的記錄
// func (l *lru) Del(key string) {
// 	if e, ok := l.cache[key]; ok {
// 		l.removeElement(e)
// 	}
// }

// // DelOldest 從 cache 中刪除最舊的記錄
// func (l *lru) DelOldest() {
// 	l.removeElement(l.ll.Front())
// }

// // 返回當前 cache 記錄數
// func (l *lru) Len() int {
// 	return l.ll.Len()
// }

// func (l *lru) removeElement(e *list.Element) {
// 	if e == nil {
// 		return
// 	}

// 	l.ll.Remove(e)
// 	en := e.Value.(*entry)
// 	l.usedBytes -= en.Len()
// 	delete(l.cache, en.key)

// 	if l.onEvicted != nil {
// 		l.onEvicted(en.key, en.value)
// 	}
// }
