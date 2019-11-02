package image

import (
	"bufio"
	"github.com/nfnt/resize"
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

// Replace替换掉原来图层的图片
func (i *Image) Replace(img *Image) *Image {
	draw.Draw(i.Img, img.Img.Bounds(), img.Img, image.ZP, draw.Src)
	return i
}

// Over覆盖在原来的图片上面
// img所选的图片
// x0,y0 其实坐标
func (i *Image) Over(img *Image, x0, y0 int) *Image {
	// 坐标
	p := image.Pt(x0, y0)
	draw.Draw(i.Img, img.Img.Bounds().Add(p), img.Img, image.ZP, draw.Over)
	return i
}

// 给图片加水印,
// src:源 可以是 Text(文字水印)
// 也可以是 Image(图片水印) ,
// x0, y0 坐标,像素点
func (i *Image) WaterMark(src interface{}, x0, y0 int) *Image {
	switch t := src.(type) {
	case Image:
		return i.Over(&t, x0, y0)
	case Text:
		return i.textMark(&t, x0, y0)
	default:
		return nil
	}
}

// 在图片上画
func (i *Image) Draw(src image.Image, sp image.Point) {
	draw.Draw(i.Img, i.Img.Bounds(), src, sp, draw.Over)
}

// textMark给图片添加文字水印
// 参数
func (i *Image) textMark(t *Text, x0, y0 int) *Image {
	t.Create()
	t.Ctx.SetDst(i.Img)
	t.Ctx.SetClip(i.Img.Bounds())

	log.Printf("%#v", *t)
	if err := t.Draw(x0, y0); err != nil {
		log.Println("水印失败", err)
	}

	return i
}

// Crop 图片裁剪成指定大小
// x0 和 y0 分别为目标图片的七点
// x1 和 y1 分别为截取的长和宽
func (i *Image) Crop(x0, y0, x1, y1 int) *Image {
	i.Img = image.NewRGBA(image.Rect(0, 0, x1, y1))
	draw.Draw(i.Img, i.Source.Bounds().Add(image.Pt(-x0, -y0)), i.Source, i.Source.Bounds().Min, draw.Src)
	return i
}

// Thumb生成缩略图
// width 和 height 分别为目标图片的长和宽
func (i *Image) Thumb(width, height uint) *Image {
	if width == 0 || height == 0 {
		width = uint(i.Source.Bounds().Max.X)
		height = uint(i.Source.Bounds().Max.Y)
	}

	i.Img = resize.Thumbnail(width, height, i.Source, resize.Lanczos3)
	return i
}

// SaveTo保存为jpg文件
// 参数path为保存路径
// 参数quality图片质量
func (i *Image) SaveTo(path string, quality int) {
	if  quality == 0{
		quality = 100
	}

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
