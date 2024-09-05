package main

import (
	"LRU-Cache/ICache"
	"fmt"
	"time"
)

func main() {
	cache := ICache.NewCache(10)
	cache.Add("five", 5)
	cache.AddWithTTL("six", 6, 10*time.Second)
	time.Sleep(4 * time.Second)
	fmt.Println(cache.Get("five"))
	fmt.Println(cache.Get("six"))
	cache.Cap()
	cache.Len()
	cache.Add("five", 5)
	cache.Remove("five")
	fmt.Println(cache.Get("five"))
	fmt.Println(cache.Get("six"))

}
