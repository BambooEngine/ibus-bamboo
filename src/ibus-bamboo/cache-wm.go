/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
