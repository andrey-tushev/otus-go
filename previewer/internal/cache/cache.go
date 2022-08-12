package cache

import (
	"os"

	"github.com/andrey-tushev/otus-go/previewer/internal/image"
)

type Cache struct {
	dir string
}

func New() Cache {
	return Cache{
		dir: "cache",
	}
}

func (c *Cache) Get(img image.PreviewImage) []byte {
	content, err := os.ReadFile(c.dir + "/" + img.HashName())
	if err != nil {
		return nil
	}

	return content
}

func (c *Cache) Set(img image.PreviewImage, content []byte) {
	f, _ := os.Create(c.dir + "/" + img.HashName())
	_, _ = f.Write(content)
	f.Close()
}
