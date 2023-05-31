package cache

import (
	"sync"
	"time"
)

// Cache 对外的缓存对象，内部维护一个接口对象ICache，主要实现ICache接口有：
// 1. LRUCache: 有淘汰机制的缓存
// 2. LRUWithTTLCache: 有淘汰机制加上过期时间的缓存
// 3. TTLCache: 有过期时间的缓存
// 且内部也维护需要传入的CacheConfig对象
type Cache struct {
	// Cache 缓存接口对象
	Cache ICache
	lock  sync.Mutex
	// Config 缓存配置项
	Config *CacheConfig
}

func NewCache(cache ICache, config *CacheConfig) *Cache {
	return &Cache{Cache: cache, Config: config, lock: sync.Mutex{}}
}

type CacheConfig struct {
	// TTL 过期时间，如果有使用，可以设置，不使用可以为空。
	// 如果需要使用，但没有设置，会默认使用10s过期时间
	TTL time.Duration
	// MaxEntries 最大缓存数
	MaxEntries int
	// Callbacks 当缓存出现修改时，可执行的回调方法
	Callbacks ChangeCallbackHandler
}

func NewCacheConfig(TTL time.Duration, maxEntries int, callbacks ChangeCallbackHandler) *CacheConfig {
	return &CacheConfig{TTL: TTL, MaxEntries: maxEntries, Callbacks: callbacks}
}

// ChangeCallbackHandler 回调接口，可提供用户实现相应方法
type ChangeCallbackHandler interface {
	OnAdd()
	OnGet()
	OnRemove()
}

// ChangeCallbackFunc 回调方法
type ChangeCallbackFunc struct {
	// OnAdd 加入缓存时，可执行的回调
	AddFunc func()
	// OnGet 获取缓存时，可执行的回调
	GetFunc func()
	// OnRemove 删除缓存时，可执行的回调
	RemoveFunc func()
}

func (c ChangeCallbackFunc) OnAdd() {
	if c.AddFunc != nil {
		c.AddFunc()
	}
}

func (c ChangeCallbackFunc) OnGet() {
	if c.GetFunc != nil {
		c.GetFunc()
	}
}

func (c ChangeCallbackFunc) OnRemove() {
	if c.RemoveFunc != nil {
		c.RemoveFunc()
	}
}

type Key interface{}

// entry 存入缓存的Value对象
type entry struct {
	key   Key
	ttl   time.Time
	value interface{}
}

const (
	maxDuration     time.Duration = 1<<63 - 1
	defaultDuration time.Duration = 10
)

// LRUCacheMode LRUCache缓存模式
func (cc *CacheConfig) LRUCacheMode() ICache {
	c := newLRU(cc.MaxEntries, maxDuration)
	return c
}

// LRUWithTTLCacheMode LRUWithTTL缓存模式，如果没有设置，就使用默认过期时间
func (cc *CacheConfig) LRUWithTTLCacheMode() ICache {
	if cc.TTL == 0 {
		cc.TTL = defaultDuration
	}
	c := newLRU(cc.MaxEntries, cc.TTL)
	return c
}

// TTLCacheMode TTL缓存模式，如果没有设置，就使用默认过期时间
func (cc *CacheConfig) TTLCacheMode(updateAgeOnGet bool) ICache {
	if cc.TTL == 0 {
		cc.TTL = defaultDuration
	}
	c := newTTLCache(cc.MaxEntries, cc.TTL, updateAgeOnGet)
	return c
}

// Add 放入缓存，如果OnAdd回调有值，就会调用
func (c *Cache) Add(key Key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Cache.add(key, value)
	if c.Config.Callbacks.OnAdd != nil {
		c.Config.Callbacks.OnAdd()
	}
}

func (c *Cache) Size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Cache.size()
}

// Get 获取缓存，如果OnGet回调有值，就会调用
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Config.Callbacks.OnGet != nil {
		c.Config.Callbacks.OnGet()
	}

	return c.Cache.get(key)
}

// Remove 删除缓存，如果OnRemove回调有值，就会调用
func (c *Cache) Remove(key Key) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Cache.remove(key)
	if c.Config.Callbacks.OnRemove != nil {
		c.Config.Callbacks.OnRemove()
	}
}

func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Cache.clear()
}
