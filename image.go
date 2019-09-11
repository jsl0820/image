package image

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

//实例化
func New(src string) *Image {
	var img Image
	f, err := os.Open(src)
	if err != nil {
		log.Println("打开文件出错:", err)
		os.Exit(-1)
	}

	config, coding, err := image.DecodeConfig(f)
	if err != nil {
		log.Println("解析图片出错%v", err)
		return nil
	}
	f.Close()

	f, err = os.Open(src)
	if err != nil {
		log.Println("打开文件出错:", err)
		os.Exit(-1)
	}

	img.Width = uint(config.Width)
	img.Height = uint(config.Height)
	img.CodingType = coding
	m, s, e := image.Decode(f)
	if err != nil {
		log.Println(e)
		return nil
	}

	img.Img = m

	log.Println(s)
	return &img
}

type Image struct {
	Src        string
	Width      uint
	Height     uint
	CodingType string
	Img        image.Image
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

//重新设置图片大小
func (i *Image) Thumb(width uint) {
	i.Copy()
	i.NewImg := resize.Resize(width, 0, i.Img, resize.Lanczos3)
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
