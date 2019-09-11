package image

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	i := New("th.jpg")
	b := i.Img.Bounds()
	log.Println("%#v", b)
}
