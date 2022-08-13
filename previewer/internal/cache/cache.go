package cache

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/andrey-tushev/otus-go/previewer/internal/preview"
)

type Cache struct {
	dir string
}

func New() Cache {
	return Cache{
		dir: "cache",
	}
}

func (c *Cache) Get(img preview.Image) *preview.Container {
	f, err := os.Open(c.filename(img))
	if err != nil {
		return nil
	}

	container := preview.NewContainer()
	dataDecoder := gob.NewDecoder(f)
	err = dataDecoder.Decode(&container)

	return container
}

func (c *Cache) Set(img preview.Image, container *preview.Container) {
	f, _ := os.Create(c.filename(img))

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(container)

	_, _ = f.Write(buff.Bytes())
	f.Close()
}

func (c *Cache) filename(img preview.Image) string {
	return c.dir + "/" + img.Key() + ".gob"
}
