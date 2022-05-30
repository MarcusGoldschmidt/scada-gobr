package main

import (
	"fmt"
)

type T = struct{}

func main() {

	ch := make(chan string, 2)
	ch <- "one"
	ch <- "two"
	close(ch)

	for elem := range ch {
		fmt.Println(elem)
	}
}
