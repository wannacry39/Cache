package ICache

import (
	"container/list"
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
	items map[any]CacheItem
	list  *list.List
	mu    sync.RWMutex
}

type CacheItem struct {
	ListElem *list.Element
	ExpireAt time.Time
	HasTTL   bool
}

func NewCache(cap int) *Cache {
	cache := &Cache{
		cap:   cap,
		items: make(map[any]CacheItem),
		list:  list.New(),
		mu:    sync.RWMutex{},
	}
	go Cleaning(cache)
	return cache
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
	c.items = make(map[any]CacheItem)
	c.list.Init()
}

func (c *Cache) AddWithTTL(key, value any, ttl time.Duration) {
	newItem := CacheItem{
		ListElem: nil,
		ExpireAt: time.Now().Add(ttl),
		HasTTL:   true,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.MoveToFront(val.ListElem)
		val.ListElem.Value = value
		newItem.ListElem = val.ListElem
		c.items[key] = newItem
		return
	}
	if c.list.Len() >= c.cap {
		backElem := c.list.Back()
		for k, v := range c.items {
			if v.ListElem == backElem {
				delete(c.items, k)
				c.list.Remove(v.ListElem)
				break
			}
		}
	}
	newItem.ListElem = c.list.PushFront(value)
	c.items[key] = newItem
}

func (c *Cache) Add(key, value any) {
	newItem := CacheItem{
		ListElem: nil,
		ExpireAt: time.Time{},
		HasTTL:   false,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.MoveToFront(val.ListElem)
		val.ListElem.Value = value
		return
	}
	if c.list.Len() >= c.cap {
		backElem := c.list.Back()
		for k, v := range c.items {
			if v.ListElem == backElem {
				delete(c.items, k)
				c.list.Remove(v.ListElem)
				break
			}
		}
	}
	newItem.ListElem = c.list.PushFront(value)
	c.items[key] = newItem
}

func (c *Cache) Get(key any) (value any, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		if time.Now().Before(val.ExpireAt) || !val.HasTTL { // true при условии, если ttl еще не прошло либо если флаг элемента HasTTL выставлен в false
			c.list.MoveToFront(val.ListElem)
			return val.ListElem.Value, true
		}

	}
	return nil, false
}

func (c *Cache) Remove(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.items[key]; ok {
		c.list.Remove(val.ListElem)
		delete(c.items, key)
	}
}

func Cleaning(c *Cache) { //функция очистки кэша от истекших значений
	ticker := time.NewTicker(time.Minute)
	ExpiredItems := []any{}
	for {
		<-ticker.C
		c.mu.RLock()
		for k, v := range c.items {
			if time.Now().After(v.ExpireAt) && v.HasTTL { // true при условии, что и ttl уже прошло и флаг HasTTL выставлен в true
				ExpiredItems = append(ExpiredItems, k)
			}
		}
		c.mu.RUnlock()
		for _, ItemKey := range ExpiredItems {
			c.Remove(ItemKey)
		}

	}
}
