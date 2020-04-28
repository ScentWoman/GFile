package gfile

import (
	"log"
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
	if path == "" {
		path = "/"
	}
	// log.Println("cache get:", path, password)
	ok = true
	list = nil

	c.mux.RLock()
	rec := c.mem[path]
	switch rec {
	case nil:
	default:
		if time.Since(rec.timestamp) > c.expire {
			log.Println("cache expired:", path)
			c.mux.RUnlock()
			c.mux.Lock()
			c.count--
			delete(c.mem, path)
			c.mux.Unlock()
			return
		}

		if password == rec.password {
			// log.Println("cache: hit")
			list = rec.list
		} else {
			log.Println("cache: wrong password")
			ok = false
		}
	}

	c.mux.RUnlock()
	return
}

func (c *cache) getWithoutPass(path string) (list []zfile.File, ok bool) {
	if path == "" {
		path = "/"
	}
	// log.Println("cache get:", path, password)
	ok = true
	list = nil

	c.mux.RLock()
	rec := c.mem[path]
	switch rec {
	case nil:
	default:
		if time.Since(rec.timestamp) > c.expire {
			log.Println("cache expired:", path)
			c.mux.RUnlock()
			c.mux.Lock()
			c.count--
			delete(c.mem, path)
			c.mux.Unlock()
			return
		}

		list = rec.list
	}

	c.mux.RUnlock()
	return
}

func (c *cache) set(path, password string, list []zfile.File) {
	// log.Println("cache: set", path, "->", list, "(pass="+password+")")
	c.mux.Lock()
	if c.count == c.size {
		for k := range c.mem {
			delete(c.mem, k)
			break
		}
	} else {
		c.count++
	}

	c.mem[path] = &record{
		timestamp: time.Now(),
		password:  password,
		list:      list,
	}
	c.mux.Unlock()
}
