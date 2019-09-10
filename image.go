package image

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func New(src string) *Image {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Printf("编码错误%v", err)
		os.Exit(-1)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return &Image{
		Src:    src,
		Width:  width,
		Height: height,
		Img:    img,
	}
}

type Image struct {
	Src    string
	Width  int
	Height int
	Img    image.Image
	NewImg *image.NRGBA
}

func (i *Image) Copy() {
	i.NewImg = image.NewNRGBA(i.Img.Bounds())
	//画
	for y := 0; y < i.Img.Bounds().Dy(); y++ {
		for x := 0; x < i.Img.Bounds().Dx(); x++ {
			i.NewImg.Set(x, y, i.Img.At(x, y))
		}
	}
}

func (i *Image) WaterMark(t Text, X, Y int) *Image {
	t.Init()
	i.Copy()
	t.Ctx.SetDst(i.NewImg)
	t.Ctx.SetClip(i.Img.Bounds())
	if err := t.Draw(X, Y); err != nil {
		log.Println(err)
	}

	return i
}

func (i *Image) SaveTo(path string, quality int) {
	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	defer out.Close()
	b := bufio.NewWriter(out)
	if err := jpeg.Encode(b, i.NewImg, &jpeg.Options{Quality: quality}); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if err := b.Flush(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
