package core

import (
	"fmt"
	"image"
)

type CardFaceId int

type CardFace struct {
	Id            CardFaceId
	SheetLocation *image.Rectangle
}

func (self *CardFace) String() string {
	return fmt.Sprintf("ID: %v, Bounds: %v", self.Id, self.SheetLocation)
	//return fmt.Sprintf("ID: %v, Bounds: (%v X %v)", self.Id, self.SheetLocation.Dx(), self.SheetLocation.Dy())
}


