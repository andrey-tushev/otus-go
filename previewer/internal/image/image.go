package image

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type PreviewImage struct {
	Path   string
	Width  int
	Height int
}

func (i PreviewImage) HashName() string {
	a := md5.Sum([]byte(i.Path))
	h := hex.EncodeToString(a[0:len(a)])
	return fmt.Sprintf("%dx%d-%s.jpg", i.Width, i.Height, h)
}
