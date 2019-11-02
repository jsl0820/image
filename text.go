package image

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
)

const (
	FONT_DPI  = 72
	FONT_SIZE = 12
	FONT_FILE = "msyh.ttf"
)

type Text struct {
	FontDPI  float64
	FontByte []byte
	FontFile string
	FontSize float64
	Color    color.RGBA
	Space    float64
	Content  []string
	Ctx      *freetype.Context
}

func (t *Text) Create() error {

	t.Ctx = freetype.NewContext()
	//设置字体文件
	if t.FontFile == "" {
		t.FontFile = FONT_FILE
	}

	// 读取字体文件如果没有设置则读取默认的字体文件
	if bt, err := ioutil.ReadFile(t.FontFile); err != nil {
		return err
	} else {
		if f, err := freetype.ParseFont(bt); err != nil {
			return err
		} else {
			t.Ctx.SetFont(f)
		}
	}

	// 如果没有DPI则设置默认的
	if t.FontDPI == 0 {
		t.FontDPI = FONT_DPI
	}

	// 设置字体大小
	if t.FontSize == 0 {
		t.FontSize = FONT_SIZE
	}

	// 设置字体间隔
	if t.Space == 0 {
		t.Space = 1.4
	}

	//默认字体颜色
	if t.Color == (color.RGBA{}) {
		t.Color = color.RGBA{0, 0, 0, 1}
	}


	t.Ctx.SetDPI(t.FontDPI)
	t.Ctx.SetFontSize(t.FontSize)
	t.Ctx.SetSrc(image.NewUniform(t.Color))
	return nil
}

// 绘制目标
func (t *Text) Draw(x0, y0 int) error {
	// 字体转换为像素点数
	h := t.Ctx.PointToFixed(t.FontSize * t.Space) >> 6
	for _, s := range t.Content {
		pt := freetype.Pt(x0, y0)
		if _, err := t.Ctx.DrawString(s, pt); err != nil {
			log.Println(err)
			return err
		}
		y0 += int(h)
	}

	return nil
}

// 绘制单个字符
//func (t *Text) DrawLetter