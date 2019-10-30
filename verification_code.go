package image

import (
	"image/color"
	"math/rand"
	"time"
	image2 "image"
)

const (
	NUMBERS          = "1234567890"
	LETTERS          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_       CodeType = iota
	NUM
	CHARS
	MIX
)

type Verify struct {
	Type     uint8
	Path     string
	Content  []byte
	Expire   uint
	UseImgBg bool
	UseNoise bool
	UseCurve bool
	FontSize uint
	FontColor []byte
	Width    uint
	Height   uint
	Length   uint8
	BgColor string
	img 	image2.Image
}

// 宽度
func (v *Verify) setWidth(){
	if v.Width == 0 {
		v.Width = 100;
	}
}

//设置高度
func (v *Verify) setHeight(){
	if v.Height == 0 {
		v.Width = v.FontSize * 10
	}
}

//设置背景颜色
func (v *Verify)setBgColor(){
	if v.BgColor == ""{
		v.BgColor = "#0000"
	}

	rgba := image2.NewRGBA(image2.Rect(0, 0, v.Width, v.Height))
	for x := 0; x < v.Width; x++  {
		for y := 0; y < v.Height; y++  {
			rgba.Set(x, y, c)
		}
	}
}


// 验证码的内容
func (v *Verify) setContent() {
	var s string
	if v.Type == NUM {
		s = NUMBERS
	}

	if v.Type == CHARS {
		s = LETTERS
	}

	if v.Type == MIX {
		s = LETTERS + NUMBERS
	}

	bytes := []byte(s)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < v.Length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	v.Content = result
}

// 生成验证码
func (v *Verify) Create() {
	v.setContent()
}

//噪点
func (v *Verify) writeNoise(){
	if v.UseNoise {
		for i := 0; i<50; i++ {
			x := rand.Intn(int(v.Width))
			for j := 0 ; j < 50 ; j++  {
				y := rand.Intn(int(v.Height))
				v := v.randColor()
				v.img.Set(X, Y, C)
			}
		}
	}
}

// 随机颜色
func (v *Verify) randColor() color.RGBA {

}


// 干扰曲线
func (v *Verify) writeCurve(){
	if v.UseCurve{

	}
}

// 检验
func (v *Verify) Check(input string) bool {
	if string(v.Content) == input {
		return true
	}

	return false
}

