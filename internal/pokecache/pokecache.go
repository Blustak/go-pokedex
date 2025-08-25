package pokecache

import (
	"sync"
	"time"
)

type Pokecache struct {
    cache map[string]cacheEntry
    mu *sync.RWMutex
}

type cacheEntry struct {
    data []byte
    createdAt time.Time
}

func NewCache(clearInterval time.Duration) Pokecache {
    cache := make(map[string]cacheEntry)
    output := Pokecache{
        cache: cache,
        mu: &sync.RWMutex{},
    }
    go output.reapLoop(clearInterval)
    return output
}

func (c *Pokecache) Add(url string, res []byte) {
    c.mu.Lock()
    c.cache[url] = cacheEntry{
        data: res,
        createdAt: time.Now(),
    }
    defer c.mu.Unlock()
}

func (c *Pokecache) Get(url string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    res, ok := c.cache[url]
    if !ok {
        return nil, false
    }
    return res.data, true
}


func (c *Pokecache) reapLoop(clearInterval time.Duration) {
    timer := time.NewTicker(clearInterval)
    var nowTime time.Time
    for {
        nowTime = <- timer.C
        c.mu.Lock()
        for k, v := range c.cache {
            if nowTime.After(v.createdAt.Add(clearInterval)) {
                delete(c.cache, k)
            }
        }
        c.mu.Unlock()
    }
}
