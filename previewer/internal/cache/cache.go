package cache

import (
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

func (c *Cache) Get(img preview.Image) []byte {
	content, err := os.ReadFile(c.dir + "/" + img.Key())
	if err != nil {
		return nil
	}

	return content
}

func (c *Cache) Set(img preview.Image, content []byte) {
	f, _ := os.Create(c.dir + "/" + img.Key())
	_, _ = f.Write(content)
	f.Close()
}
