package image

import (
	"bufio"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
)

//判断图片编码类型
func EncodeType(src string) string {
	f, err := os.Open(src)
	if err != nil {
		log.Println(err)
	}

	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		log.Println(err)
	}

	return http.DetectContentType(buff)
}

//实例化
func New(src string) *Image {
	f, err := os.Open(src)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	defer f.Close()
	var img image.Image

	if EncodeType(src) == "image/jpeg" {
		img, err = jpeg.Decode(f)
		if err != nil {
			log.Printf("jpeg图片解码出错%v", err)
		}
	}

	if EncodeType(src) == "image/png" {
		img, err = png.Decode(f)
		if err != nil {
			log.Printf("png图片解码出错%v", err)
		}
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

//拷贝一张图片到
func (i *Image) Copy() {
	i.NewImg = image.NewNRGBA(i.Img.Bounds())
	//画
	for y := 0; y < i.Img.Bounds().Dy(); y++ {
		for x := 0; x < i.Img.Bounds().Dx(); x++ {
			i.NewImg.Set(x, y, i.Img.At(x, y))
		}
	}
}

// 给图片加水印,
// src:源 可以是 Text(文字水印)
// 也可以是 Image(图片水印) ,
// X, Y 坐标,像素点
func (i *Image) WaterMark(src interface{}, X, Y int) *Image {
	switch t := src.(type) {
	case Image:
		return i.imgMark(t, X, Y)
	case Text:
		return i.textMark(t, X, Y)
	default:
		return nil
	}
}

//图片水印
func (i *Image) imgMark(img Image, X, Y int) *Image {
	i.Copy()
	off := image.Pt(X, Y)

	draw.Draw(i.NewImg, i.Img.Bounds(), i.Img, image.ZP, draw.Src)
	draw.Draw(i.NewImg, img.Img.Bounds().Add(off), img.Img, image.ZP, draw.Over)

	return i
}

//文字水印
func (i *Image) textMark(t Text, X, Y int) *Image {
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
