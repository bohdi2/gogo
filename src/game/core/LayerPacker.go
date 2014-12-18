package core

import (
	"fmt"
	"image"
	"log"
	"sort"
)

var nextCardFaceId = make(chan CardFaceId)

type Info struct {
	id    CardFaceId
	width int
	filename string
}

type InfoList []Info

func newInfo(width int, filename string) Info {
	return Info{<- nextCardFaceId, width, filename}
}

func (infoList InfoList) Len() int {
	return len(infoList)
}

func (infoList InfoList) Swap(i, j int) {
    infoList[i], infoList[j] = infoList[j], infoList[i]
}

func (infoList InfoList) Less(i, j int) bool {
	return infoList[i].width > infoList[j].width
}

type layerPacker struct {
	filenames []string
	images    []image.Image
}

func NewLayerPacker() SheetPacker {
	return &layerPacker{make([]string, 0), make([]image.Image, 0)}
}

func (packer *layerPacker) Add(filename string) {
	packer.filenames = append(packer.filenames, filename)
}

// Writes
func (packer layerPacker) Pack(dirname, outputName string) {
	offset := image.Point{0, 0}
	columnWidth := 0
	sheetWidth := 0
	sheetHeight := 0
	pairs := make([]*pair, 0, 100)
	infoList := make(InfoList, 0, 100)

	sheetOutputName := dirname + "/" + outputName + ".jpg"
	dataOutputName := dirname + "/" + outputName + ".json"

	fmt.Printf("OutputNames %v, %v\n", sheetOutputName, dataOutputName)

	for _, filename := range packer.filenames {
		config, _, err := readImageConfig(filename)
		if err != nil {
			log.Printf("Error reading %v, Err %v\n", filename, err)
			continue
		}
		info := newInfo(config.Width, filename)
		fmt.Println(info)
		infoList = append(infoList, info)
	}

	sort.Sort(infoList)

	for _, info := range infoList {

		picture, err := readImage(info.filename)
		if err != nil {
			log.Printf("Error reading %v, Err %v\n", info.filename, err)
			continue
		}

		bounds := enlarge(5, picture.Bounds())

		sheetLocation := image.Rect(offset.X, offset.Y, offset.X+bounds.Dx(), offset.Y+bounds.Dy())

		pairs = append(pairs, &pair{picture, CardFace{info.id, &sheetLocation}})

		if sheetLocation.Dx() > columnWidth {
			columnWidth = sheetLocation.Dx()
		}

		offset.Y += sheetLocation.Dy()
		if sheetHeight < offset.Y {
			sheetHeight = offset.Y
		}

		if offset.Y > 10000 {
			sheetWidth += columnWidth
			columnWidth = 0

			offset = image.Point{sheetWidth, 0}
		}
	}

	fmt.Printf("1 Writing Sheet Image\n");
	if err := writeSheetImage(sheetOutputName, sheetWidth+columnWidth, sheetHeight, pairs); err != nil {
		log.Printf("Error writing sheet %v, Err %v\n", sheetOutputName, err)
	}

	fmt.Printf("2 Writing CardFace Data\n")
	if err := writeCardFaceData(pairs); err != nil {
		log.Printf("Error writing CardFace Data %v, Err %v\n", dataOutputName, err)
	}

	fmt.Printf("3\n")
}

func init() {
	go func() {
		for i := 0; ; i++ {
			nextCardFaceId <- CardFaceId(i)
		}
	}()
}
