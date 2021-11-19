package cache_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/iscreen/go-cache"
	"github.com/iscreen/go-cache/lru"
)

func BenchmarkTourCacheSetParallel(b *testing.B) {
	cache := cache.NewTourCache(nil, lru.New(b.N*100, nil))
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0

		for pb.Next() {
			cache.Set(paralleKey(id, counter), value())
			counter++
		}
	})
}

// func key(i int) string {
// 	return fmt.Sprintf("key-%04d", i)
// }

func value() []byte {
	return make([]byte, 100)
}

func paralleKey(threadID, counter int) string {
	return fmt.Sprintf("key-%4d-%06d", threadID, counter)
}
