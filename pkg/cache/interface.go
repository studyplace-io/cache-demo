package cache

// ICache 接口对象，需要实现缓存各个方法
type ICache interface {
	// add 放入缓存
	add(key Key, value interface{})
	// size 数量
	size() int
	// get 获取缓存
	get(key Key) (value interface{}, ok bool)
	// remove 删除缓存
	remove(key Key)
	// clear 清理所有缓存
	clear()
}
