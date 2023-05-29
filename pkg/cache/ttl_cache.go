package cache

import (
	"sort"
	"time"
)

// ttlCache 实现ttl过期时间的缓存，使用
type ttlCache struct {
	maxEntries     int
	updateAgeOnGet bool
	cache          map[Key]*entry
	// 使用expiration维护
	expiration map[int64]Key
	expiry     time.Duration
}

func newTTLCache(maxEntries int, expiry time.Duration, updateAgeOnGet bool) *ttlCache {
	return &ttlCache{
		maxEntries:     maxEntries,
		updateAgeOnGet: updateAgeOnGet,
		cache:          make(map[Key]*entry),
		expiration:     make(map[int64]Key),
		expiry:         expiry,
	}
}

func (c *ttlCache) add(key Key, value interface{}) {
	// 1. 如果没有就创建map
	if c.cache == nil {
		c.cache = make(map[Key]*entry)
		c.expiration = make(map[int64]Key)
	}
	// 2. 如果本来就有，就放入map中，并更新ttl
	if ee, ok := c.cache[key]; ok {
		ee.value = value
		c.changeTTL(key)
		return
	}

	// 3. 创建 entry，并放入map中
	ele := &entry{
		ttl:   time.Now().Add(c.expiry),
		value: value,
		key:   key,
	}
	c.cache[key] = ele
	exp := ele.ttl.UnixNano()
	c.expiration[exp] = key

	// 4. 如果长度超过，必须删除最后一个
	if c.maxEntries != 0 && len(c.cache) > c.maxEntries {
		c.purgeToCapacity()
	}
}

// changeTTL 更新ttl
func (c *ttlCache) changeTTL(key Key) {
	if ee, ok := c.cache[key]; ok {
		delete(c.expiration, ee.ttl.UnixNano())
		ee.ttl = time.Now().Add(c.expiry)
		exp := ee.ttl.UnixNano()
		c.expiration[exp] = key
	}
}

func (c *ttlCache) size() int {
	if c.cache == nil {
		return 0
	}
	return len(c.cache)
}

// get 获取缓存
func (c *ttlCache) get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	// 如果获取到，先查看是否过期，如果过期直接返回，
	// 如果需要更新ttl，则更新
	if ele, hit := c.cache[key]; hit {
		if time.Now().After(ele.ttl) {
			c.remove(key)
			return
		}
		if c.updateAgeOnGet {
			c.changeTTL(key)
		}
		return ele.value, hit
	}
	return
}

// remove 删除缓存
func (c *ttlCache) remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *ttlCache) purgeToCapacity() {
	expKeys := make([]int64, 0, len(c.expiration))
	for k := range c.expiration {
		expKeys = append(expKeys, k)
	}
	// 存小到大排序
	sort.Slice(expKeys, func(i, j int) bool { return expKeys[i] < expKeys[j] })
	for _, k := range expKeys {
		if len(c.cache) <= c.maxEntries && k > time.Now().UnixNano() {
			return
		} else {
			c.remove(c.expiration[k])
		}
	}
}

// removeElement 删除元素
func (c *ttlCache) removeElement(e *entry) {
	delete(c.expiration, e.ttl.UnixNano())
	delete(c.cache, e.key)
}

// clear 清理
func (c *ttlCache) clear() {
	c.cache = nil
	c.expiration = nil
}
