package cache

import (
	"container/list"
	"time"
)

// lruCache 实现LRU淘汰机制的缓存，使用链表记录所有对象的value
type lruCache struct {
	maxEntries int
	ll         *list.List
	cache      map[Key]*list.Element
	expiry     time.Duration
}

func newLRU(maxEntries int, expiry time.Duration) *lruCache {
	return &lruCache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
		expiry:     expiry,
	}
}

// add 加入缓存
func (c *lruCache) add(key Key, value interface{}) {
	// 1. 如果没有map，先创建map
	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.ll = list.New()
	}
	// 2. 如果能从map中找到，先放到链表最前面，更新 ttl 与 value
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).ttl = time.Now().Add(c.expiry)
		ee.Value.(*entry).value = value
		return
	}
	// 3. 创建 entry 放入最链表前端，并放入map中
	ele := c.ll.PushFront(&entry{key, time.Now().Add(c.expiry), value})
	c.cache[key] = ele
	// 4. 如果长度超过，必须删除最后一个
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *lruCache) size() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// get 获取缓存
func (c *lruCache) get(key Key) (value interface{}, ok bool) {

	if c.cache == nil {
		return
	}

	// 如果获取到，先查看是否过期，如果过期直接返回，
	// 没有过期就放入链表前头，并返回
	if ele, hit := c.cache[key]; hit {
		if time.Now().After(ele.Value.(*entry).ttl) {
			c.remove(key)
			return
		}
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// remove 删除缓存
func (c *lruCache) remove(key Key) {
	if c.cache == nil {
		return
	}
	// 找到就删除
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// removeOldest 删除最老的
func (c *lruCache) removeOldest() {
	if c.cache == nil {
		return
	}
	// 找到最老的ele，并删除
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

// removeElement 删除元素
func (c *lruCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// clear 清除链表与
func (c *lruCache) clear() {
	c.ll = nil
	c.cache = nil
}
