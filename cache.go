package gfile

import (
	"sync"
	"time"

	"github.com/ScentWoman/GFile/zfile"
)

type cache struct {
	mux         sync.RWMutex
	count, size int
	mem         map[string]*record
	expire      time.Duration
}

type record struct {
	timestamp time.Time
	password  string
	list      []zfile.File
}

func newCache(size int, expire time.Duration) *cache {
	c := &cache{
		mux:    sync.RWMutex{},
		size:   size,
		mem:    make(map[string]*record),
		expire: expire,
	}
	go func() {
		c.autoDelete()
	}()
	return c
}

func (c *cache) autoDelete() {
	for {
		time.Sleep(24 * time.Hour)
		c.mux.Lock()
		for k, v := range c.mem {
			if time.Since(v.timestamp) > c.expire {
				c.count--
				delete(c.mem, k)
			}
		}
		c.mux.Unlock()
	}
}

func (c *cache) get(path, password string) (list []zfile.File, ok bool) {
	ok = true
	list = nil

	c.mux.RLock()
	rec := c.mem[path]
	switch rec {
	case nil:
	default:
		if time.Since(rec.timestamp) > c.expire {
			c.mux.RUnlock()
			c.mux.Lock()
			c.count--
			delete(c.mem, path)
			c.mux.Unlock()
			return
		}

		if password == rec.password {
			list = rec.list
		} else {
			ok = false
		}
	}

	c.mux.RUnlock()
	return
}

func (c *cache) set(path, password string, list []zfile.File) {
	c.mux.Lock()
	if c.count != c.size {
		c.count++
		c.mem[path] = &record{
			timestamp: time.Now(),
			password:  password,
			list:      list,
		}
	}
	c.mux.Unlock()
}
