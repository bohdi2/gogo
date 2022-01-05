gogo
====

This is a game I wrote ten plus years ago in golang. The idea is simple: you have a pile of cards and you have to
remove the topmost cards. The tricky part is I wanted to have piles in excess of a million cards.
I never got much higher than 100,000 cards because I kept everything in memory.

At large numbers of cards the trivial/obvious graph algorithms that I knew did not scale well. The core problem were

* Determining whether a mouse click was on top of one of the root nodes (a top most card).
* Updating the set of root nodes when a card was removed.
* Redrawing parts of a card when it became partially exposed.
* Maintaining a transitive reduction of the graph of cards.
