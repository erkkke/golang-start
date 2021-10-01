package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	recursiveSend(t, ch)
	close(ch)
}

func recursiveSend(t * tree.Tree, ch chan int) {
	if t != nil {
		recursiveSend(t.Left, ch)
		ch <- t.Value
		recursiveSend(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		first, ok1 := <-ch1
		second, ok2 := <-ch2
		if ok1 != ok2 || first != second {
			return false
		}
		if (!ok1) {
			break
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
