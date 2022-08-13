package cache

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"

	"github.com/andrey-tushev/otus-go/previewer/internal/preview"
)

type Cache struct {
	dir string
	mu  sync.RWMutex
}

func New() Cache {
	return Cache{
		dir: "cache",
	}
}

func (c *Cache) Get(img preview.Image) *preview.Container {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f, err := os.Open(c.filename(img))
	if err != nil {
		return nil
	}
	defer f.Close()

	container := preview.NewContainer()
	dataDecoder := gob.NewDecoder(f)
	err = dataDecoder.Decode(&container)

	return container
}

func (c *Cache) Set(img preview.Image, container *preview.Container) {
	c.mu.Lock()
	defer c.mu.Unlock()

	f, err := os.Create(c.filename(img))
	if err != nil {
		return
	}
	defer f.Close()

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(container)

	_, _ = f.Write(buff.Bytes())
}

func (c *Cache) filename(img preview.Image) string {
	return c.dir + "/" + img.Key() + ".gob"
}
