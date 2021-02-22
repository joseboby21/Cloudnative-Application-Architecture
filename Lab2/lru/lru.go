package lru

import (
	"errors"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string,0)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	// Your code here....
	k,ok := key.(string)
	if !ok{
		return nil,errors.New("Not a valid input")
	}
	if _, ok := lru.cache[k]; !ok {
		return nil, errors.New("Not present in Cache")
	}
	lru.qDel(k)
	lru.queue = append(lru.queue, k)
	return lru.cache[k], nil
}

func (lru *lruCache) Put(key, val interface{}) error {
	// Your code here....
	k,ok1 := key.(string)
	
	if !ok1{
		return errors.New("Not a valid input")
	}

	v,ok2 := val.(string)

	if !ok2{
		return errors.New("Not a valid input")
	}

	_,ok3 := lru.cache[k]

	if ok3{
		lru.cache[k] = v
		lru.qDel(k)
		lru.queue = append(lru.queue, k)
	}else{
		if lru.remaining == 0 {
			delete(lru.cache, lru.queue[0])
			lru.qDel(lru.queue[0])
			lru.remaining++
		}
		lru.cache[k] = v
		lru.queue = append(lru.queue, k)
		lru.remaining--
	}	
	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
