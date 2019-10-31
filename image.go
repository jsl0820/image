package image

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

//实例化
func New(src string) *Image {

	f, err := os.Open(src)
	if err != nil {
		log.Println("打开文件出错:", err)
		os.Exit(-1)
	}
	var img *Image
	if config, decode, err := image.DecodeConfig(f); err != nil {
		log.Println("获取图片信息失败:", err)
		return nil
	} else {
		img = &Image{
			Width:      uint(config.Width),
			Height:     uint(config.Height),
			DecodeType: decode,
		}
	}

	// 必须限关掉重新打开才正常运行
	f.Close()
	f, err = os.Open(src)
	if source, _, err := image.Decode(f); err != nil {
		log.Println("解析图片出错:", err)
		return nil
	} else {
		img.Create()
		img.Source = source
		draw.Draw(img.Img, source.Bounds(), source, image.ZP, draw.Src)
	}

	return img

}

type Image struct {
	BgColor    color.RGBA
	Width      uint
	Height     uint
	DecodeType string
	Source     image.Image
	Img        draw.Image
}

// 生成带背景色的画布
func (i *Image) Blank() *Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{i.BgColor}, image.ZP, draw.Src)
	i.Img = img
	return i
}

// 生成空白的画布
func (i *Image) Create() *Image {
	i.Img = image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))
	return i
}

// 替换掉原来图层的图片
func (i *Image) Replace(img *Image) *Image {
	draw.Draw(i.Img, img.Img.Bounds(), img.Img, image.ZP, draw.Src)
	return i
}

// 覆盖在原来的图片上面
func (i *Image) Over(img *Image, x0, y0 int) *Image {
	// 坐标
	p := image.Pt(x0, y0)
	draw.Draw(i.Img, img.Img.Bounds().Add(p), img.Img, image.ZP, draw.Over)
	return i
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
func (i *Image) imgMark(img Image, x0, y0 int) *Image {
	i.Create()
	off := image.Pt(x0, y0)
	draw.Draw(i.Img, i.Img.Bounds(), i.Img, image.ZP, draw.Src)
	draw.Draw(i.Img, img.Img.Bounds().Add(off), img.Img, image.ZP, draw.Over)

	return i
}

// 在图片上画
func (i *Image) Draw(src image.Image, sp image.Point) {
	draw.Draw(i.Img, i.Img.Bounds(), src, sp, draw.Over)
}

// 给图片添加文字水印
// 参数
func (i *Image) textMark(t Text, X, Y int) *Image {
	t.Init()
	i.Create()
	t.Ctx.SetDst(i.Img)
	t.Ctx.SetClip(i.Img.Bounds())
	if err := t.Draw(X, Y); err != nil {
		log.Println(err)
	}

	return i
}


// 保存为jpg文件
// 参数path为保存路径
// 参数quality图片质量
func (i *Image) SaveTo(path string, quality int) {
	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	defer out.Close()
	b := bufio.NewWriter(out)
	if err := jpeg.Encode(b, i.Img, &jpeg.Options{Quality: quality}); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if err := b.Flush(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
