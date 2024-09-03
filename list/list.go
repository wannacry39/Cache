package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type ICache interface {
	Cap() int
	Len() int
	Clear() // удаляет все ключи
	Add(key, value any)
	AddWithTTL(key, value any, ttl time.Duration) // добавляет ключ со сроком жизни ttl
	Get(key any) (value any, ok bool)
	Remove(key any)
}

type Cache struct {
	cap   int
	items map[any]*list.Element
	list  *list.List
	mu    sync.RWMutex
}

func WaitForExpire(c *Cache, key any, ttl time.Duration) {
	select {
	case <-time.After(ttl):
		c.Remove(key)
		return
	}
}

func NewCache(cap int) *Cache {
	return &Cache{
		cap:   cap,
		items: make(map[any]*list.Element),
		list:  list.New(),
		mu:    sync.RWMutex{},
	}
}

func (c *Cache) Cap() int {
	return c.cap
}

func (c *Cache) Len() int {
	c.mu.RLock()
	length := c.list.Len()
	c.mu.RUnlock()
	return length
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[any]*list.Element)
	c.list.Init()
}

func (c *Cache) Add(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.MoveToFront(val)
		val.Value = value
		return
	}
	if c.list.Len() >= c.cap {
		backElem := c.list.Back()
		for k, v := range c.items {
			if v == backElem {
				delete(c.items, k)
				c.list.Remove(v)
				break
			}
		}

	}
	elem := c.list.PushFront(value)
	c.items[key] = elem
}

func (c *Cache) AddWithTTL(key any, value any, ttl time.Duration) {
	c.Add(key, value)
	go WaitForExpire(c, key, ttl)

}

func (c *Cache) Get(key any) (value any, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

func (c *Cache) Remove(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.Remove(val)
		delete(c.items, key)
	}
}

func main() {
	cache := NewCache(3)
	go func() {
		cache.AddWithTTL("one", 1, 17*time.Second)
		cache.AddWithTTL("two", 2, 16*time.Second)
		cache.AddWithTTL("three", 3, 1*time.Second)
		cache.Add("two", 6)
	}()

	go func() {
		cache.Get("one")
		cache.Get("two")
		cache.Get("three")
	}()

	time.Sleep(2 * time.Second)
	for k, v := range cache.items {
		fmt.Printf("k: %d, v: %d ", k, v.Value)
	}
	fmt.Println()

	for e := cache.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%d ", e.Value)
	}

}
