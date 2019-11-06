package image

import (
	"image/color"
	"log"
	"testing"
)

func TestVerifyCreate(t *testing.T) {
	verify := Verify{
		Type:    MIX,
		Width:   240,
		Height:  80,
		Length:  6,
		UseNoise: true,
		UseCurve: true,
		BgColor: color.RGBA{225, 255, 255, 255},
	}

	verify.Create()
	log.Printf("验证码字符串%#v", verify.Content)
}
