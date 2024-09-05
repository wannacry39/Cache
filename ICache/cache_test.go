package ICache

import (
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	cache := NewCache(10)

	for i := 0; i < 100; i++ {
		cache.Add(i, i)
		if cache.Len() > cache.Cap() {
			t.Errorf("length greater then capacity")
			break
		}

		if len(cache.items) != cache.list.Len() {
			t.Errorf("map and list length's are not equal")
			break
		}
	}
}

func TestAdd(t *testing.T) {
	cache := NewCache(5)
	for i := 1; i < 11; i++ {
		cache.Add(i, i)
		if i != cache.list.Front().Value {
			t.Errorf("Added element should be first")
			break
		}

		val, ok := cache.items[i]
		if !ok {
			t.Errorf("Key not added in map")
			break
		}
		if val.ListElem.Value != i {
			t.Errorf("values in map and list are not equal")
			break
		}
	}
	cache.Add(10, "newvalue")
	val := cache.items[10]
	if val.ListElem.Value != "newvalue" {
		t.Errorf("existing key with an old value")
	}
}

func TestAddWithTTL(t *testing.T) {
	cache := NewCache(10)
	for i := 0; i < 5; i++ {
		cache.Add(i, i)
	}
	cache.AddWithTTL("one", 1, 300*time.Millisecond)
	if cache.list.Front().Value != 1 {
		t.Errorf("Added value should be first in the list")
		return
	}
	cache.AddWithTTL("two", 2, 500*time.Millisecond)
	if cache.list.Front().Value != 2 {
		t.Errorf("Added value should be first in the list")
		return
	}

	cache.AddWithTTL("two", 22, 1*time.Minute)
	val := cache.items["two"]
	if val.ListElem.Value != 22 {
		t.Errorf("existing key with an old value")
	}
}

func TestGet(t *testing.T) {
	cache := NewCache(10)
	for i := 1; i < 11; i++ {
		cache.Add(i, i)
	}

	for i := 1; i < 11; i++ {
		val, ok := cache.Get(i)
		if val != i {
			t.Errorf("incorrect value")
			break
		}
		if !ok {
			t.Errorf("bool flag of existing value is false")
			break
		}
		if val != cache.list.Front().Value {
			t.Errorf("requested value should be in the front of the list")
			break
		}
	}

	val, ok := cache.Get("AnyKey")
	if val != nil {
		t.Errorf("Not existing value should equals nil")
		return
	}
	if ok != false {
		t.Errorf("bool flag of not existing value is true")
		return
	}

	cache.AddWithTTL("one", 1, 300*time.Millisecond)
	cache.AddWithTTL("two", 2, 500*time.Millisecond)

	time.Sleep(300 * time.Millisecond)

	_, Exists1 := cache.Get("one")
	_, Exists2 := cache.Get("two")

	if Exists1 != false {
		t.Errorf("Got a value after it's expired")
		return
	}
	if Exists2 != true {
		t.Errorf("didn't get a value before it's expired")
		return
	}

}

func TestClear(t *testing.T) {
	cache := NewCache(100)
	for i := 0; i < 100; i++ {
		cache.Add(i, i)
	}

	cache.Clear()

	if len(cache.items) != 0 {
		t.Errorf("map's length should be 0")
		return
	}
	if cache.list.Len() != 0 {
		t.Errorf("list's length should be 0")
	}

}

func TestRemove(t *testing.T) {
	cache := NewCache(5)
	cache.Add("key", "value")

	cache.Remove("key")

	val, ok := cache.items["key"]
	if ok != false {
		t.Errorf("removed key in map")
		return
	}
	if val.ListElem != nil {
		t.Errorf("removed value in list")
	}
}
