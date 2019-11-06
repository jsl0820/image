package image

import (
	"image/color"
	"log"
	"testing"
)


func TestNew(t *testing.T) {
	img := New("aaa.jpg")
	log.Printf("%v", img)
}

func TestSaveTo(t *testing.T) {
	img := New("aaa.jpg")
	img.SaveTo("save.jpg", 10)
}

func TestBlank(t *testing.T) {

	img := &Image{
		Width:      100,
		Height:     100,
		BgColor:    color.RGBA{255, 0, 0, 1},
	}

	img.Blank().SaveTo("save.jpg", 10)
}

func TestCreate(t *testing.T) {
	img := &Image{
		Width:      100,
		Height:     100,
	}

	img.Blank().SaveTo("save.jpg", 10)
}

func TestOver(t *testing.T) {
	img1 := New("th.jpg")
	img2 := New("aaa.jpg")

	img2.Over(img1, 10, 10).SaveTo("save.jpg", 10)
}

func TestReplace(t *testing.T) {
	img1 := &Image{
		Width:      400,
		Height:     400,
	}

	img2 := New("th.jpg")
	img1.Create().Replace(img2).SaveTo("save.jpg", 10)
}

func TestWaterMark(t *testing.T) {
	var text = Text{
		Content:[]string{"ERTYU"},
		Color:color.RGBA{255,0, 0, 255},
	}

	var img = New("aaa.jpg")
	img.WaterMark(text, 100, 100 ).SaveTo("save.jpg", 100)
}

//
func TestImage_Crop(t *testing.T) {
	img := New("th.jpg")
	img.Crop(20, 20, 250, 250).SaveTo("save.jpg", 10)
}

func TestImage_Thumb(t *testing.T) {
	img := New("testdata/th.jpg")
	img.Thumb(10, 10).SaveTo("testdata/thhumb.jpeg", 10)
}

