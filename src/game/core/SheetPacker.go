package core

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
        "launchpad.net/gobson/bson"
        "launchpad.net/mgo"
	"log"
	"os"
)

type SheetPacker interface {
	Add(filename string)
	Pack(dirname, filename string)
}

type pair struct {
	CardFaceImage image.Image
	CardFace      CardFace
}

type tuple struct {
	Id     CardFaceId
        Location image.Rectangle
	CardFace image.Image
}

func newTuple(p pair) tuple {
	return tuple{p.CardFace.Id, *p.CardFace.SheetLocation, p.CardFaceImage}
}

func (p *pair) height() int {
	return p.CardFace.SheetLocation.Dy()
}

func (p *pair) width() int {
	return p.CardFace.SheetLocation.Dx()
}

func shrink(n int, r image.Rectangle) image.Rectangle {
	p := image.Pt(n,n)
	return image.Rectangle{r.Min.Add(p), r.Max.Sub(p)}
}

func enlarge(n int, r image.Rectangle) image.Rectangle {
	p := image.Pt(n,n)
	return image.Rectangle{r.Min.Sub(p), r.Max.Add(p)}
}




func readImageConfig(filename string) (image.Config, string, os.Error) {
	source, err := os.Open(filename)
	if err != nil {
		return image.Config{}, "", err
	}
	defer source.Close()

	return image.DecodeConfig(source)
}

func readImage(filename string) (image.Image, os.Error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	return jpeg.Decode(source)
}

// Get the width x height data for the face

func readPair(filename string, id CardFaceId, offset image.Point) (*pair, os.Error) {
	faceImage, err := readImage(filename)
	if err != nil {
		return nil, err
	}

	bounds := enlarge(5, faceImage.Bounds())
	sheetLocation := image.Rect(offset.X, offset.Y, offset.X+bounds.Dx(), offset.Y+bounds.Dy())

	return &pair{faceImage, CardFace{id, &sheetLocation}}, nil
}

func writeCardFaceData(pairs []*pair) os.Error {
        session, err := mgo.Mongo("127.0.0.1")
        if err != nil {
                panic(err)
        }
        defer session.Close()

	fmt.Printf("Got Mongo\n")

        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        c := session.DB("gogo").C("faces")
	fmt.Printf("Got Mongo faces\n")

	c.RemoveAll(bson.M{})
	fmt.Printf("Removed faces\n")

	for _, pair := range pairs {
		fmt.Printf("  Add pair %v\n", pair)
		err = c.Insert(newTuple(*pair))
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func writeSheetImage(filename string, width, height int, pairs []*pair) os.Error {
	// Create the face strip image.

	fmt.Printf("Creating new image (%v x %v) %v\n", width, height, width*height)
	sheet := image.NewRGBA(width, height)

	// Once more iterate over the faces, but this time
	// copy them into the face sheet.

	for _, pair := range pairs {
		draw.Draw(sheet, *pair.CardFace.SheetLocation, image.Black, image.Point{0, 0}, draw.Over)
		draw.Draw(sheet, shrink(5, *pair.CardFace.SheetLocation), pair.CardFaceImage, image.Point{5, 5}, draw.Over)
	}

	fmt.Printf("Saving face sheet as %v\n", filename)

	// Save the face sheet.
	dest, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer dest.Close()
	fmt.Printf("returning encode\n")

	return jpeg.Encode(dest, sheet, nil)
}


func readCardFaceData() *[]CardFace {
        session, err := mgo.Mongo("127.0.0.1")
        if err != nil {
                panic(err)
        }
        defer session.Close()


        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        collection := session.DB("gogo").C("faces")

	faces := make([]CardFace, 0)
	var result *tuple

	err = collection.Find(nil).For(&result, func() os.Error {
		faces = append(faces, CardFace{result.Id, &result.Location})
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("CardFaces[]: %v\n", faces)
	return &faces
}
