package lru

import "container/list"

type Cache struct {
	maxBytes    int64
	nbytes      int64
	ll          *list.List
	cache       map[string]*list.Element
	OnEvicted   func(key string, value Value)
}

type entry struct {
	key    string
	value  Value
}

type Value interface {
	Len()   int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, val Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.value.Len())
		kv.value = val
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: val})
		c.nbytes += int64(val.Len()) + int64(len(key))
		c.cache[key] = ele
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.removeOldest()
	}
}

func (c *Cache) Get(key string) (val Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		val = ele.Value.(*entry).value
		c.ll.MoveToFront(ele)
		return val, true
	}
	return
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		c.nbytes -= int64(kv.value.Len()) + int64(len(kv.key))
		delete(c.cache, kv.key)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}