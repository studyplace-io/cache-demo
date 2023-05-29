package test

import (
	"fmt"
	"github.com/cache-example/pkg/cache"
	"log"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {

	config := cache.NewCacheConfig(0, 20, cache.ChangeCallbacks{
		OnAdd: func() {
			fmt.Println("entry add...")
		},
		OnGet: func() {
			fmt.Println("entry get...")
		},
		OnRemove: func() {
			fmt.Println("entry delete...")
		},
	})

	lruCache := cache.NewCache(config.LRUCacheMode(), config)

	lruCache.Add(1, "this is test 1")
	lruCache.Add(2, "this is test 2")
	lruCache.Add(3, "this is test 3")
	lruCache.Add(4, "this is test 4")
	lruCache.Add(5, "this is test 5")
	lruCache.Add(6, "this is test 6")

	log.Println("init LRU lruCache finish")
	log.Printf("lruCache size: %d", lruCache.Size())

	if v, ok := lruCache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	if v, ok := lruCache.Get(3); ok {
		log.Printf("get key (%d) success, value: %s", 3, v)
	}

	lruCache.Add(7, "this is test 7")
	lruCache.Add(8, "this is test 8")

	if _, ok := lruCache.Get(4); !ok {
		log.Printf("get key (%d) failed", 4)
	}
}

func TestLRUWithTTLCache(t *testing.T) {
	config := cache.NewCacheConfig(4, 20, cache.ChangeCallbacks{
		OnAdd: func() {
			fmt.Println("entry add...")
		},
		OnGet: func() {
			fmt.Println("entry get...")
		},
		OnRemove: func() {
			fmt.Println("entry delete...")
		},
	})

	lruWithTTLCache := cache.NewCache(config.LRUWithTTLCacheMode(), config)

	lruWithTTLCache.Add(1, "this is test 1")
	lruWithTTLCache.Add(2, "this is test 2")
	lruWithTTLCache.Add(3, "this is test 3")
	lruWithTTLCache.Add(4, "this is test 4")
	lruWithTTLCache.Add(5, "this is test 5")
	lruWithTTLCache.Add(6, "this is test 6")

	log.Println("init LRU TTL cache finish")
	log.Printf("cache size: %d", lruWithTTLCache.Size())

	if v, ok := lruWithTTLCache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
	}

	time.Sleep(5 * time.Second)

	if _, ok := lruWithTTLCache.Get(2); !ok {
		log.Printf("get key (%d) failed", 2)
	}
}

func TestTTLCache(t *testing.T) {
	config := cache.NewCacheConfig(time.Duration(10), 20, cache.ChangeCallbacks{
		OnAdd: func() {
			fmt.Println("entry add...")
		},
		OnGet: func() {
			fmt.Println("entry get...")
		},
		OnRemove: func() {
			fmt.Println("entry delete...")
		},
	})

	tllCache := cache.NewCache(config.TTLCacheMode(true), config)

	tllCache.Add(1, "this is test 1")
	time.Sleep(3 * time.Millisecond)
	tllCache.Add(2, "this is test 2")
	time.Sleep(3 * time.Millisecond)
	tllCache.Add(3, "this is test 3")
	time.Sleep(3 * time.Millisecond)
	tllCache.Add(4, "this is test 4")
	time.Sleep(3 * time.Millisecond)
	tllCache.Add(5, "this is test 5")
	time.Sleep(3 * time.Millisecond)
	tllCache.Add(6, "this is test 6")
	time.Sleep(3 * time.Millisecond)

	log.Println("init TTL cache finish")
	log.Printf("cache size: %d", tllCache.Size())

	time.Sleep(3 * time.Second)

	if v, ok := tllCache.Get(2); ok {
		log.Printf("get key (%d) success, value: %s", 2, v)
		log.Printf("cache size: %d", tllCache.Size())
	}

	time.Sleep(3 * time.Second)

	if _, ok := tllCache.Get(3); !ok {
		log.Printf("get key (%d) failed", 3)
		log.Printf("cache size: %d", tllCache.Size())
	}

	log.Println("Add new three item")
	tllCache.Add(7, "this is test 7")
	tllCache.Add(8, "this is test 8")
	tllCache.Add(9, "this is test 9")
	log.Printf("cache size: %d", tllCache.Size())

	if v, ok := tllCache.Get(8); ok {
		log.Printf("get key (%d) success, value: %s", 8, v)
	}

}

