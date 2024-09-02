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
	return c.list.Len()
}

func (c *Cache) Clear() {
	c.items = make(map[any]*list.Element)
	c.list.Init()
}

func (c *Cache) Add(key, value any) {
	if val, ok := c.items[key]; ok {
		c.items[key] = c.list.InsertAfter(value, val)
		c.list.Remove(c.items[key].Prev())
		return
	}
	if c.Cap() <= c.Len() {
		for k, v := range c.items {
			if v == c.list.Back() {
				delete(c.items, k)
				c.list.Remove(v)
				c.items[key] = c.list.PushBack(value)
				return
			}
		}
	}

	c.items[key] = c.list.PushBack(value)
	return
}

func (c *Cache) Get(key any) (value any, ok bool) {
	if val, ok := c.items[key]; ok {
		data := c.list.Remove(val)
		c.list.PushFront(data)
		return data, true
	}
	return nil, false
}

func (c *Cache) Remove(key any) {
	if _, ok := c.items[key]; !ok {
		return
	}
	c.list.Remove(c.items[key])
	delete(c.items, key)
}

func main() {

	cache := NewCache(3)

	cache.Add("one", 1)
	cache.Add("two", 2)
	cache.Add("three", 3)
	cache.Add("four", 4)

	for k, v := range cache.items {
		fmt.Printf("k: %s, v: %d ", k, v.Value)
	}
	fmt.Println()

	for e := cache.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%d ", e.Value)
	}

	cache.Get("two")
	cache.Remove("one")

	fmt.Println()
	for k, v := range cache.items {
		fmt.Printf("k: %s, v: %d ", k, v.Value)
	}
	fmt.Println()

	for e := cache.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%d ", e.Value)
	}

}
