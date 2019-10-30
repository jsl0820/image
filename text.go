package image

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
)


type Text struct {
	DPI      float64
	FontByte []byte
	FontSrc  string
	Size     float64
	Color    color.RGBA
	Space    float64
	Content  []string
	Ctx      *freetype.Context
}

func (t *Text) LoadFontFile() *Text {
	b, err := ioutil.ReadFile(t.FontSrc)
	if err != nil {
		log.Println(err)
	}

	t.FontByte = b
	return t
}

func (t *Text) Init() *Text {
	t.LoadFontFile()
	f, err := freetype.ParseFont(t.FontByte)
	if err != nil {
		log.Println(err)
	}

	t.Ctx = freetype.NewContext()
	t.Ctx.SetDPI(t.DPI)
	t.Ctx.SetFont(f)
	t.Ctx.SetFontSize(t.Size)
	t.Ctx.SetSrc(image.NewUniform(t.Color))
	return t
}

//
func (t *Text) Draw(X, Y int) error {
	h := t.Ctx.PointToFixed(t.Size*t.Space) >> 6
	for _, s := range t.Content {
		pt := freetype.Pt(X, Y)
		if _, err := t.Ctx.DrawString(s, pt); err != nil {
			log.Println(err)
			return err
		}
		Y += int(h)
	}

	return nil
}
