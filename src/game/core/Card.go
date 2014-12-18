package core

import (
	"image"
	"strconv"
)

type CardId int

func (self CardId) IsBackground() bool {
	return self < 0
}

func (self CardId) String() string {
	return strconv.Itoa(int(self))
}


type Card struct {
	CardId   CardId
	Face     *image.Rectangle
	Location *image.Rectangle
	Reverse  map[CardId][]*image.Rectangle
	//VisableFragments []*image.Rectangle       // Visable fragments. Initially one big fragment.
}

func NewCard(cardId CardId, face *image.Rectangle, x1, y1, x2, y2 int) Card {
	r := image.Rect(x1, y1, x2, y2)
	return Card{cardId,
		face,
		&r,
		make(map[CardId][]*image.Rectangle, 4),
		//make([]*image.Rectangle, 1),
	}
}


func (self *Card) Overlaps(card Card) bool {
	return self.Location.Overlaps(*card.Location)
}

func (self Card) ContainsPoint(x, y int) bool {
	return image.Point{x, y}.In(*self.Location)
}

type DrawCommand struct {
	Type        string
	Width       int
	Height      int
	FacePoint image.Point
	GroundPoint image.Point
}

type ClearCommand struct {
	Type string
}


func (self *Card) Command(fragment *image.Rectangle) DrawCommand {
	return DrawCommand{"DRAW",
		fragment.Dx(),
		fragment.Dy(),
		fragment.Min.Sub(self.Location.Min).Add(self.Face.Min),
		fragment.Min}

}
