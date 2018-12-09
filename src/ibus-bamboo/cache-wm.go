package main

import (
	"sync"
)

type CacheWM struct {
	sync.Mutex
	values   map[uint32][]string
	keyList  []uint32
	maxItems int
}

func NewCacheWM(maxItems int) *CacheWM {
	return &CacheWM{
		values:   map[uint32][]string{},
		maxItems: maxItems,
	}
}

func (c *CacheWM) Get(window uint32) ([]string, bool) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.values[window]
	if ok {
		lenKeyList := len(c.keyList)
		for i := lenKeyList - 1; i >= 0; i-- {
			if c.keyList[i] == window {
				for j := i + 1; j < lenKeyList; j++ {
					c.keyList[j-1] = c.keyList[j]
				}
				c.keyList[lenKeyList-1] = window
				break
			}
		}
	}

	return v, ok
}

func (c *CacheWM) Set(window uint32, wm []string) {
	c.Lock()
	defer c.Unlock()
	c.values[window] = wm
	if len(c.keyList) < c.maxItems {
		c.keyList = append(c.keyList, window)
	} else {
		delete(c.values, c.keyList[0])
		c.keyList = append(c.keyList[1:], window)
	}
}
