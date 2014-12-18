package core

import (
	"fmt"
	"gob"
	"image"
	"os"
	"rand"
)

func foo(s []*image.Rectangle, x1, y1, x2, y2 int) []*image.Rectangle {
	if x1 != x2 && y1 != y2 {
		r := image.Rect(x1, y1, x2, y2)
		return append(s, &r)
	}
	return s

}

// top intersects bottom, bottom is split into small rectangles.

func split(top, bottom *image.Rectangle) ([]*image.Rectangle, *image.Rectangle) {
	s := make([]*image.Rectangle, 0, 4)

	overlap := top.Intersect(*bottom)

	//fmt.Printf("Fragment %v, is being split by overlap %v\n", fragment, overlap)

	s = foo(s, bottom.Min.X, bottom.Min.Y, bottom.Max.X, overlap.Min.Y)
	s = foo(s, bottom.Min.X, overlap.Min.Y, overlap.Min.X, overlap.Max.Y)
	s = foo(s, overlap.Max.X, overlap.Min.Y, bottom.Max.X, overlap.Max.Y)
	s = foo(s, bottom.Min.X, overlap.Max.Y, bottom.Max.X, bottom.Max.Y)

	return s, &overlap
}


type Pile struct {
	Width            int
	Height           int
	CardFaces        []CardFace                    // Card images
	NextId           CardId
	Cards            map[CardId]*Card
	VisableCards     map[CardId]int                // Set of visable cards
	Roots            map[CardId]int                // Set of root cards
	VisableFragments map[CardId][]*image.Rectangle // Visable fragments. Initially one big fragment.
}

type PileBuilder struct {
	Width            int
	Height           int
	NextId           CardId
	CardFaces        []CardFace                    // Card images
}

func (self *PileBuilder) nextId() CardId {
	id := self.NextId
	self.NextId++
	return id
}

func NewPileBuilder(size, width, height int) *PileBuilder {
	fmt.Printf("NewPile(%d, %d, %d)\n", size, width, height)

	p := new(PileBuilder)
	p.Width = width
	p.Height = height

	p.CardFaces = *readCardFaceData()

	return p
}

func (self *Pile) nextId() CardId {
	id := self.NextId
	self.NextId++
	return id
}

func NewPile(size, width, height int) *Pile {
	fmt.Printf("NewPile(%d, %d, %d)\n", size, width, height)

	p := new(Pile)
	p.Width = width
	p.Height = height

	p.CardFaces = *readCardFaceData()
	for _, face := range p.CardFaces {
		fmt.Printf("RawCardFace %v\n", face)
	}

	p.VisableFragments = make(map[CardId][]*image.Rectangle, 10)
	p.Roots = make(map[CardId]int, 10)
	p.Cards = make(map[CardId]*Card, 10)
	p.VisableCards = make(map[CardId]int, 10)

	return p
}

func (self *Pile) CreateBackground() {
	id := CardId(-1)

	faceId := CardFaceId(0)
	r := self.CardFaces[faceId].SheetLocation

	x := 0
	y := 0

	for x < self.Width {
		y = 0
		for y < self.Height {
			self.Add(NewCard(id, r, x, y, x+r.Dx(), y+r.Dy()))
			id--
			y += r.Dy()
		}

		x += r.Dx()
	}
}

func (self *Pile) CreateCard() Card {
	id := self.nextId()

	i := rand.Intn(len(self.CardFaces)-1) + 1 // Don't include faceId 0
	faceId := CardFaceId(i)

	r := self.CardFaces[faceId].SheetLocation
	x1 := rand.Intn(self.Width - r.Dx())
	y1 := rand.Intn(self.Height - r.Dy())
	x2 := x1 + r.Dx()
	y2 := y1 + r.Dy()

	return NewCard(id, r, x1, y1, x2, y2)
}

func (self *Pile) getOverlappedCards(card Card) []CardId {
	result := make([]CardId, 0, len(self.VisableCards))

	for child, _ := range self.VisableCards {
		if card.Overlaps(*self.Cards[child]) {
			//fmt.Printf("GOC: %v, %v\n", card.Location, self.Cards[child].Location)
			result = append(result, child)
		}
	}
	return result
}

func (self *Pile) Add(card Card) {
	id := card.CardId

	self.Cards[id] = &card

	for rootId, _ := range self.Roots {
		if card.Overlaps(*self.Cards[rootId]) {
			self.Roots[rootId] = 0, false
		}
	}

	self.Roots[id] = 0

	for _, overlappedId := range self.getOverlappedCards(card) {
		visableFragments := self.VisableFragments[overlappedId]
		newVisableFragments := make([]*image.Rectangle, 0, len(visableFragments))

		for _, visableFragment := range visableFragments {
			if visableFragment.Overlaps(*card.Location) {
				visable, hidden := split(card.Location, visableFragment)
				newVisableFragments = append(newVisableFragments, visable...)
				card.Reverse[overlappedId] = append(card.Reverse[overlappedId], hidden)
			} else {
				newVisableFragments = append(newVisableFragments, visableFragment)
			}
		}

		if len(newVisableFragments) > 0 {
			self.VisableFragments[overlappedId] = newVisableFragments
		} else {
			self.VisableFragments[overlappedId] = newVisableFragments, false
			self.VisableCards[overlappedId] = 0, false
		}

	}

	self.VisableCards[id] = 0
	self.VisableFragments[id] = make([]*image.Rectangle, 0)
	self.VisableFragments[id] = append(self.VisableFragments[id], card.Location)
}

func (self *Pile) isUncovered(id CardId) bool {
	card := self.Cards[id]
	area := card.Location.Dx() * card.Location.Dy()

	for _, visableFragment := range self.VisableFragments[id] {
		area -= visableFragment.Dx() * visableFragment.Dy()
	}

	return 0 == area
}

func (self *Pile) findRoot(x, y int) (*Card, os.Error) {
	for rootId, _ := range self.Roots {
		if self.Cards[rootId].ContainsPoint(x, y) {
			return self.Cards[rootId], nil
		}
	}
	return nil, os.NewError("core: Root not found.")
}

func (self *Pile) Remove(x, y int) (map[CardId][]*image.Rectangle, os.Error) {

	var err os.Error
	if card, err := self.findRoot(x, y); err == nil {
		self.VisableFragments[card.CardId] = nil, false

		for cardId, fragments := range card.Reverse {
			if cardId.IsBackground() {
				continue
			}
			//fmt.Printf("Remove %s, %d, %v\n", cardId, len(self.VisableFragments[cardId]), self.VisableFragments[cardId])
			self.VisableFragments[cardId] = append(self.VisableFragments[cardId], fragments...)
			//fmt.Printf("           %d, %v\n", len(self.VisableFragments[cardId]), self.VisableFragments[cardId])

			if int(cardId) >= 0 && self.isUncovered(cardId) {
				self.Roots[cardId] = 0
			}

		}

		// put roots back, put fragments back
		self.Cards[card.CardId] = nil, false
		self.Roots[card.CardId] = 0, false
		self.VisableFragments[card.CardId] = nil, false
		return card.Reverse, nil
	}
	return nil, err
}

func (self Pile) String() string {
	var s string
	s += fmt.Sprintf("Area: (%d x %d), ", self.Width, self.Height)
	s += fmt.Sprintf("CardFaces: [%d]. ", len(self.CardFaces))
	s += fmt.Sprintf("Visable: %v, ", len(self.VisableCards))
	s += fmt.Sprintf("Roots: %v, ", len(self.Roots))

	count1 := 0
	for _, card := range self.Cards {
		for _, visableFragments := range card.Reverse {
			count1 += len(visableFragments)
		}
	}

	s += fmt.Sprintf("Cards: %v, %d ", len(self.Cards), count1)

	count2 := 0
	for _, visableFragments := range self.VisableFragments {
		count2 += len(visableFragments)
	}

	s += fmt.Sprintf("Fragments: %v, %d %d", len(self.VisableFragments), count2, count1 + count2)

	return s
}

func (self *Pile) Visit(f func(pile *Pile, card *Card)) {
	for _, x := range self.Cards {
		f(self, x)
	}
}


func (self *Pile) VisitFragments(f func(pile *Pile, id CardId, fragments []*image.Rectangle)) {
	fmt.Printf("VF %d\n", len(self.VisableFragments))
	for cardId, visableFragments := range self.VisableFragments {
		f(self, cardId, visableFragments)
	}
}


func (self *Pile) Store(fname string) os.Error {
	fmt.Printf("Saving %v\n", fname)

	f, err := os.Create(fname)
	if err != nil {
		fmt.Printf("--ERRROR-- storing to %v\n", fname)
		return err
	}

	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(self)
}

func Load(fname string) (*Pile, os.Error) {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("--ERRROR-- loading from  %v\n", fname)
		return nil, err
	}

	defer f.Close()

	var pile Pile
	decoder := gob.NewDecoder(f)
	return &pile, decoder.Decode(&pile)
}


func init() {
	//rand.Seed(time.Nanoseconds() % 1e9)
	rand.Seed(0)
}
