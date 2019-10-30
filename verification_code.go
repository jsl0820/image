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
	Content  string
	Expire   uint
	UseImgBg bool
	UseNoise bool
	UseCurve bool
	FontSize float32
	FontColor color.RGBA
	Width    uint
	Height   uint
	Length   uint8
	BgColor color.RGBA
	img 	image2.Image
}

// 宽度
func (v *Verify) setWidth() uint {
	if v.Width == 0 {
		size := float32(v.FontSize)
		length := float32(v.Length)
		v.Width = uint(length * size * 1.5 + length * size / 2)
	}

	return v.Width
}

//设置高度
func (v *Verify) setHeight() uint {
	if v.Height == 0 {
		size := float32(v.FontSize)
		v.Height =  uint(size * 2.5)
	}

	return v.Height
}


// 验证码的内容
func (v *Verify) setContent() string{
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
	for i := 0; i < int(v.Length); i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	v.Content = string(result)
	return string(result)
}

// 生成验证码
func (v *Verify) Create() {

	conts := []string{v.setContent()}
	//字体图片
	textImg := &Text{
		Color: v.FontColor,
		FontSrc: v.setContent(),
		Content: conts,
	}

	// 背景图片
	bgIamge := &Image{
		Color: v.BgColor,
		Width: v.setWidth(),
		Height:v.setHeight(),
	}

	v.img = bgIamge.Draw()
	v.writeNoise()
	v.writeNoise()

	//设置验证码图片
	//生成字体图片图片
	// 字体内容
	// 噪点
	// 干扰曲线
}


//噪点
func (v *Verify) writeNoise(){
	if v.UseNoise {
		for i := 0; i<50; i++ {
			x := rand.Intn(int(v.Width))
			for j := 0 ; j < 50 ; j++  {
				y := rand.Intn(int(v.Height))
				c := v.randColor()
				v.img.Set(x, y, c)
			}
		}
	}
}

// 随机颜色
func (v *Verify) randColor() color.RGBA {
	r :=  rand.Intn(255)
	g :=  rand.Intn(255)
	b :=  rand.Intn(255)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 1}
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

