package main

import (
	cache "LRU-Cache/Cache"
	"fmt"
	"sync"
	"time"
)

func main() {
	cache := cache.NewCache(10)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		cache.Add("one", 1)
		cache.Add("two", 2)
		cache.Add("three", 3)
		cache.AddWithTTL("four", 4, 1*time.Second)
		cache.AddWithTTL("five", 5, 6*time.Second)

	}()

	time.Sleep(5 * time.Second)

	go func() {
		fmt.Println(cache.Get("one"))
		fmt.Println(cache.Get("two"))
		fmt.Println(cache.Get("three"))
		fmt.Println(cache.Get("four"))
		fmt.Println(cache.Get("five"))
		wg.Done()
	}()

	wg.Wait()
}
