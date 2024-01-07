package go_concurrency

import (
	"sync"
	"sort"
)

type ConcurrentMap[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

type ConcurrentMapItem[K comparable, V any] struct {
	Key   K
	Value V
}

func NewMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		items: make(map[K]V),
	}
}

func (c *ConcurrentMap[K, V]) Set(key K, value V) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = value
}

func (c *ConcurrentMap[K, V]) Remove(key K) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *ConcurrentMap[K, V]) Get(key K) (v V, ok bool) {
	c.Lock()
	defer c.Unlock()

	v, ok = c.items[key]

	return v, ok
}

func (c *ConcurrentMap[K, V]) Has(key K) (ok bool) {
	c.Lock()
	defer c.Unlock()

	_, ok = c.items[key]

	return ok
}

func (c *ConcurrentMap[K, V]) Iter() <-chan ConcurrentMapItem[K, V] {
	iChan := make(chan ConcurrentMapItem[K, V])
	f := func() {
		c.Lock()
		defer c.Unlock()
		for k, v := range c.items {
			iChan <- ConcurrentMapItem[K, V]{k, v}
		}
		close(iChan)
	}
	go f()

	return iChan
}

func (c *ConcurrentMap[K, V]) Length() int {
	c.Lock()
	defer c.Unlock()
	return len(c.items)
}

func (c *ConcurrentMap[K, V]) Sort(reverse bool) {
	c.Lock()
	defer c.Unlock()

	type kv struct{
		Key K
		Value V
	}

	var tmp []kv
	for k,v := c.items {
		tmp = append(tmp, kv{Key: k, Value: v})
	}

	sort.Slice(tmp, func(i, j int) bool{
		if(reverse){
			return tmp[i].Value > tmp[j].value
		}

		return tmp[i].Value < tmp[j].value
	})

	c.items = make(map[K]V)

	for k,v := tmp{
		c.items[k] = v
	}
}