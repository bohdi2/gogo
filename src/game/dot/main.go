package main

// Dot Generator

import (
	"flag"
	"fmt"
	. "game/core"
)




func main() {
	flag.Parse()
	filename := flag.Arg(0)

	pile, err := Load(filename)
	if err != nil {
		fmt.Printf("Load error %v for %v\n", err, filename)
	}

	fmt.Println("digraph pile {")
	pile.Visit(func(pile *Pile, card *Card) {
		//fmt.Printf("Card: %v\n", card)
		//fmt.Printf("Children: %v\n", pile.Children)
		//for _,childId := range pile.Children.Row(card.CardId) {
		//	fmt.Printf("c%v ->  c%v\n", card.CardId, childId)
		//}
	})
	fmt.Println("}")

	//fmt.Printf("Pile just loaded: %v\n", pile)

}

