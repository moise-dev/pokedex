package pokecache

import "time"

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.entries[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	entity, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entity.val, true

}

func (c *Cache) replLoop(interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for k, entry := range c.entries {
					if time.Since(entry.createdAt) > interval {
						delete(c.entries, k)
					}
				}
			}
		}
	}()
}

func NewCache(duration time.Duration) Cache {
	c := Cache{
		entries: make(map[string]cacheEntry),
	}
	c.replLoop(duration)
	return c
}
