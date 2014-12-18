package main

// Web Server

import (
	"flag"
	"fmt"
	. "game/core"
	//	"image"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
)

type Event struct {
	Type string
	X, Y int
}

func createData(filename string, size, width, height int) {
	fmt.Printf("createData() %s, %d, %d, %d\n", filename, size, width, height)

	builder := NewPileBuilder(size, width, height)
	pile := NewPile(size, width, height)

	pile.CreateBackground()

	dot := size / 10

	for ii := 0; ii < size; ii++ {
		pile.Add(pile.CreateCard())
		if 0 == dot || 0 == (ii%dot) {
			fmt.Printf(".")
		}
	}

	//pile.VisitFragments(func(pile *Pile, id CardId, fragments []*image.Rectangle) {
	//	fmt.Printf("Card %s, len %d, cap %d\n", id, len(fragments), cap(fragments))
	//})

	//pile.Add(NewCard(2, 7, 109, 148, 382, 332))
	//pile.Add(NewCard(3, 8,  96, 187, 315, 417))

	err := pile.Store(filename)
	fmt.Printf("Save err %v\n", err)
}

func main() {
	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memProfile = flag.String("memprofile", "", "write memory profile to file")
	var memProfileRate = flag.Int("memprofilerate", runtime.MemProfileRate, "read godoc runtime MemProfileRate")

	flag.Parse()

	if *memProfileRate != runtime.MemProfileRate {
		runtime.MemProfileRate = *memProfileRate
	}

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatalf("can't create %s: %s", *memProfile, err)
		}
		defer func() {
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatalf("can't write %s: %s", *memProfile, err)
			}
			f.Close()
		}()
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatalf("can't create %s: %s", *cpuProfile, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if 4 != flag.NArg() {
		fmt.Printf("Not enough args %v", flag.Args())
	}

	filename := flag.Arg(0)
	size, _ := strconv.Atoi(flag.Arg(1))
	width, _ := strconv.Atoi(flag.Arg(2))
	height, _ := strconv.Atoi(flag.Arg(3))

	fmt.Printf("Args: %s, %d, (%d, %d)\n", filename, size, width, height)

	createData(filename, size, width, height)

}
