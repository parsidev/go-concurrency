package go_concurrency

import (
	"sort"
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
	less  func(a, b V, rev bool) bool
}

type ConcurrentMapItem[K comparable, V any] struct {
	Key   K
	Value V
}

func NewMapNullableSort[K comparable, V any]() ConcurrentMap[K, V] {
	return NewMap[K, V](nil)
}

func NewMap[K comparable, V any](less func(a, b V, rev bool) bool) ConcurrentMap[K, V] {
	return ConcurrentMap[K, V]{
		items: make(map[K]V),
		less:  less,
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

func (c *ConcurrentMap[K, V]) Keys() (keys []K){
	keys = make([]K, 0, len(c.items))

	for k := range c.items{
		keys = append(keys, k)
	}

	return keys
}


func (c *ConcurrentMap[K, V]) First() (item *ConcurrentMapItem[K, V]) {
	c.Lock()
	defer c.Unlock()

	for k, v := range c.items{
		item = &ConcurrentMapItem[K, V]{k, v}
		break
	}

	return item
}

func (c *ConcurrentMap[K, V]) Sort(reverse bool) {
	c.Lock()
	defer c.Unlock()

	if c.less == nil {
		return
	}

	type kv struct {
		Key   K
		Value V
	}

	var tmp []kv
	for k, v := range c.items {
		tmp = append(tmp, kv{Key: k, Value: v})
	}

	sort.Slice(tmp, func(i, j int) bool {
		return c.less(tmp[i].Value, tmp[j].Value, reverse)
	})

	c.items = make(map[K]V)

	for _, item := range tmp {
		c.items[item.Key] = item.Value
	}
}
