package image

import (
	"log"
	"testing"
)

func TestCreate(t *testing.T){
	ins := &Verify{
		Type:1,
	}

	log.Println(ins.Content)
}