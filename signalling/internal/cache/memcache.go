package cache

var MemcacheInstance = NewMemcache()

type Memcache struct {
	cache map[string]*interface{}
	dataCache map[string]*[]byte
}

func NewMemcache() *Memcache {
	return &Memcache{
		cache: make(map[string]*interface{}),
		dataCache: make(map[string]*[]byte),
	}
}

func (c *Memcache) Set(key string, value interface{}, data []byte) {
	c.cache[key] = &value
	c.dataCache[key] = &data
}

func (c *Memcache) Get(key string) (interface{}, bool) {
	value, ok := c.cache[key]
	return value, ok
}

func (c *Memcache) GetData(key string) (*[]byte, bool) {
	value, ok := c.dataCache[key]
	return value, ok
}

func (c *Memcache) All() map[string]*[]byte {
	return c.dataCache
}

func (c *Memcache) Delete(key string) {
	delete(c.cache, key)
}

func (c *Memcache) ClearData(cl interface{}) {
	for k, v := range c.cache {
		if v == cl {
			delete(c.cache, k)
			delete(c.dataCache, k)
		}
	}
}

func (c *Memcache) Contains(key string) bool {
	_, ok := c.cache[key]
	return ok
}
