package concurrent

import (
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"sync"
)

type SyncMap[Key comparable, Value any] struct {
	mx sync.RWMutex
	m  map[Key]Value
}

func (c *SyncMap[Key, Value]) Get(key Key) (Value, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *SyncMap[Key, Value]) Set(key Key, value Value) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *SyncMap[Key, Value]) Remove(key Key) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, key)
}

func (c *SyncMap[Key, Value]) Range() (values []Value) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	for _, v := range c.m {
		values = append(values, v)
	}
	return values
}

func (c *SyncMap[Key, Value]) Find(searchedItem Value, comparator func(Value, Value) bool) (*Key, bool) {
	for key, value := range c.m {
		if comparator(searchedItem, value) {
			return &key, true
		}
	}
	return nil, false
}

func (c *SyncMap[Key, Value]) TryAdd(value Value, keygen func() Key, comparator func(Value, Value) bool) (Value, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if comparator == nil {
		panic(apperrors.ComparatorNotFound)
	}
	for _, v := range c.m {
		if comparator(v, value) {
			return v, false
		}
	}
	c.m[keygen()] = value
	return value, true
}

func (c *SyncMap[Key, Value]) AddRange(values []Value, keyExtractor func(Value) Key) {
	c.mx.Lock()
	defer c.mx.Unlock()
	for _, v := range values {
		c.m[keyExtractor(v)] = v
	}
}

func NewSyncMap[Key comparable, Value any]() *SyncMap[Key, Value] {
	return &SyncMap[Key, Value]{
		m: make(map[Key]Value),
	}
}
