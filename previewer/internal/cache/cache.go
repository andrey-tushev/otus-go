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

func (c *Cache) Get(img preview.Image) *preview.Container {
	content, err := os.ReadFile(c.dir + "/" + img.Key())
	if err != nil {
		return nil
	}

	container := preview.NewContainer()
	container.Body = content

	return container
}

func (c *Cache) Set(img preview.Image, container *preview.Container) {
	f, _ := os.Create(c.dir + "/" + img.Key())
	_, _ = f.Write(container.Body)
	f.Close()
}
