package image

import (
	"image/color"
	"log"
	"testing"
)

func TestVerifyCreate(t *testing.T) {
	verify := Verify{
		Type:    MIX,
		Length:  4,
		//UseNoise: true,
		//UseCurve: true,
		CurveWeight: 3,
		BgColor: color.RGBA{225, 255, 255, 255},
	}

	verify.Create()
	log.Printf("验证码字符串%#v", verify.Content)
}
