package main

import (
	//"container/list"
	"fmt"
	I "image"
	"testing"
)


func TestClip(t *testing.T) {
	painting := new(Painting)

	fmt.Printf("Painting: %v\n", painting)

	painting.addSpriteFragment(newSpriteFragment(0, 0, 100, 100))
	fmt.Printf("Painting: %v\n", painting)

	painting.addSpriteFragment(newSpriteFragment(30, 30, 80, 70))
	fmt.Printf("Painting: %v\n", painting)
	/*
		segments := list.New()
		segments.PushFront(&r)

		fmt.Printf("R:\n")
		for e := segments.Front(); e != nil; e = e.Next() {
			c := e.Value.(*SpriteFragment);
			fmt.Printf("  %v\n", c)
	        }

		//r = I.Rect(20, 20, 30, 30)
		//segment(segments, &r)

		fmt.Printf("R:\n")
		//for e := segments.Front(); e != nil; e = e.Next() {
		//	c := e.Value.(*I.Rectangle);
		//	fmt.Printf("  %v\n", c)
	        //}
	*/
}

func TestSplit(t *testing.T) {
	segment := newSpriteFragment(0, 0, 100, 100)
	small := I.Rect(30, 30, 80, 70)

	for _, s := range segment.split(small) {
		fmt.Printf("%v\n", s)
	}
}
