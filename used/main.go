package main

import (
	"context"
	"fmt"
)

var dd chan int

func main() {
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return // returning not to leak the goroutine
			//case dst <- n:
			//	n++
		}
	}
	//c := make(chan int)
	//c <- 10
	//fmt.Println(c)
	//
	//v := <- c
	//fmt.Println(v)

	//select {
	//case c <- 10:
	//default:
	//
	//}
}

func handler() {
	fmt.Println("tewsts")
	dd <- 0
	fmt.Println("nihao")
}
