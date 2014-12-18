package main

// CreateSprites

import (
	"flag"
	//"fmt"
	. "game/core"
)


// Usage: main dest_dir dest_filename_tempalte  source1.jpg ...
//
// Concatinates one or more source jpg files into a single jpg file. The result (dest.jpg) is
// a tall skinny file containing the source images.

func main() {
	flag.Parse()

	outputDir := flag.Arg(0)
	outputFilename := flag.Arg(1)
	jpegFilenames := flag.Args()[2:]

	packer := NewLayerPacker()

	for _, filename := range jpegFilenames {
		packer.Add(filename)
	}

	packer.Pack(outputDir, outputFilename)
}
