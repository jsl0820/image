package image

import (
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"strings"
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
	CurveWeight int
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
	log.Println("宽度", v.Width)
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
	if (v.FontSize == 0){
		v.FontSize = 25
	}

	bgImg := Image{
		Width:   v.setWidth(),
		Height:  v.setHeight(),
		BgColor: v.BgColor,
	}

	log.Printf("%#v", *v)
	// 背景画布
	bgImg.Blank()

	// 验证码字符串
	strs := strings.Split(v.Content, "")
	x0 := (float32(RandInt(12, 16)) / 10 )  *  v.FontSize

	for i :=0 ; i < len(strs); i++ {
		text := Text{
			FontDPI: 200,
			FontSize: float64(v.FontSize),
			Color:    v.randFontColor(),
			Content:  []string{strs[i]},
		}

		y0 := 20
		x0 =+ x0 +  (float32(RandInt(12, 16)) / 10 )  * 1.2  *  v.FontSize
		//bgImg.WaterMark(text,  int(x0), int(y0))
		log.Println("距离", x0, y0)
		bgImg.WaterMark(text, int(x0), y0)
	}


	//噪点
	if v.UseNoise {
		for i := 0 ; i < 10;i++  {
			v.writeNoiseg(bgImg.Img)
		}
	}

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
			go func(X, Y int , C color.RGBA){
				img.Set(X, Y, C)
				img.Set(X+1, Y, C)
			}(X, Y, C);
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

// 随机字体颜色
func (v *Verify) randFontColor() color.RGBA {
	rand.Seed(time.Now().UnixNano())

	R := RandInt(1, 150)
	G := RandInt(1, 150)
	B := RandInt(1, 150)

	return color.RGBA{uint8(R), uint8(G), uint8(B), 1}
}


// 画干扰曲线
func (v *Verify) writeCurve(img draw.Image) {
	height := float64(v.Height)
	rand.Seed(time.Now().UnixNano())
	a0 := int(math.Ceil(height / 2))
	b0 := int(math.Ceil(height / 4))
	t0 := rand.Intn(int(v.Height) + int(v.Width) * 2) - int(v.Height)
	x := rand.Intn(int(v.Width))
	rgba := v.randColor()
	curve1 := &SinCurve{
		A : float64(rand.Intn(a0)),
		B : float64(rand.Intn(2 * b0) - b0) + (height / 2 ),
		W : float64(2 * math.Pi / float64(t0)),
	}

	log.Printf("%#v", curve1)

	v.setPixel(img,0, x, rgba, curve1)

	rand.Seed(time.Now().UnixNano())
	a1 := int(math.Ceil(height / 2))
	t1 := rand.Intn(int(v.Height) + int(v.Width) * 2) - int(v.Height)

	A1 := float64(RandInt(0, a1))
	W1 := float64(2 * math.Pi / float64(t1))
	B2 := curve1.DirectY(float64(x)) - (A1 * math.Sin(W1 * float64(x))) - (height / 2 )

	curve2 := &SinCurve{
		A : A1,
		B : B2 + (height / 2 ),
		W : W1,
	}

	v.setPixel(img, x, int(v.Width) , rgba, curve2)
}

type SinCurve struct {
	A float64
	W float64
	B float64
	F float64
}

func (s *SinCurve)DirectY(x float64) float64{
	return s.A * math.Sin(s.W * x + s.F) + s.B
}

// 干扰线
func (v *Verify)setPixel(img draw.Image, x0, x1 int,  color color.RGBA, curve * SinCurve)*Verify{
	for i := x0; i < x1; i++ {
		go func(i int) {
			X := float64(i)
			Y := curve.DirectY(X)
			log.Println(X, Y)
			for i := 0 ; i < v.CurveWeight ; i ++  {
				img.Set(int(math.Ceil(X)), int(math.Ceil(Y + float64(i))), color)
			}
		}(i)
	}
	return v
}

// 对比
func (v *Verify) Check(input string) bool {
	if string(v.Content) == input {
		return true
	}

	return false
}

// 两整数之间取随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max + min ) - min
}

