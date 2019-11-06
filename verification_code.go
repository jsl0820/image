package image

import (
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	NUMBERS = "1234567890"
	LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUM     = iota
	CHARS
	MIX
)

type Verify struct {
	Type      uint8
	Path      string
	Content   string
	Expire    uint
	UseImgBg  bool
	UseNoise  bool
	UseCurve  bool
	FontSize  float32
	FontColor color.RGBA
	Width     uint
	Height    uint
	Length    uint8
	BgColor   color.RGBA
}

// 宽度
func (v *Verify) setWidth() uint {
	if v.Width == 0 {
		size := float32(v.FontSize)
		length := float32(v.Length)
		v.Width = uint(length*size*1.5 + length*size/2)
	}

	return v.Width
}

//设置高度
func (v *Verify) setHeight() uint {
	if v.Height == 0 {
		size := float32(v.FontSize)
		v.Height = uint(size * 2.5)
	}

	return v.Height
}

// 验证码的内容
func (v *Verify) setContent() string {
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
func (v *Verify) Create() *Verify {
	v.setContent()
	bgImg := Image{
		Width:   v.Width,
		Height:  v.Height,
		BgColor: v.BgColor,
	}

	// 背景画布
	bgImg.Blank()

	// 验证码字符串
	text := Text{
		FontDPI: 200,
		FontSize: 16,
		Color:   color.RGBA{225, 0, 0,255},
		Content: []string{v.Content},
	}

	bgImg.WaterMark(text, 15, 60)
	// 噪点
	//if v.UseNoise {
	//	for i := 0 ; i < 10;i++  {
	//		v.writeNoiseg(bgImg.Img)
	//	}
	//}

	// 曲线
	if v.UseCurve {
		v.writeCurve(bgImg.Img)
	}

	bgImg.SaveTo("testdata/save.jpeg", 100)
	return v
}

//噪点
func (v *Verify) writeNoiseg(img draw.Image ) {
	for i := 0; i<50; i++ {
		X := rand.Intn(int(v.Width))
		Y := rand.Intn(int(v.Height))

		for j := 0 ; j < 50 ; j++  {
			C := v.randColor()
			img.Set(X, Y, C)
			img.Set(X+1, Y, C)
			img.Set(X-1, Y, C)
		}
	}
}

// 随机颜色
func (v *Verify) randColor() color.RGBA {
	R := rand.Intn(150 + 225) - 150
	G := rand.Intn(150 + 225) - 150
	B := rand.Intn(150 + 225) - 150

	return color.RGBA{uint8(R), uint8(G), uint8(B), 1}
}

// 干扰曲线
// 曲线函数：Y=Asin(WX+φ)+B
func (v *Verify) writeCurve(img draw.Image ) {
	//height := float64(v.Height)

	//A := rand.Intn(int(math.Ceil(height / 2)))
	//B := height / 4
	//W := math.Pi
	C := v.randColor()


	for i := 0; i < int(v.Width); i++ {
		X := float64(i)
		Y := 70 * math.Sin(math.Pi*X) + 40

		log.Println(X, Y)
		img.Set(int(math.Ceil(X)), int(math.Ceil(Y)), C)
	}
}



// 检验
func (v *Verify) Check(input string) bool {
	if string(v.Content) == input {
		return true
	}

	return false
}
