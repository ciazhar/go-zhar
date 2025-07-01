package main

import "fmt"

func generate() <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	for n := range square(generate()) {
		fmt.Println(n)
	}
}
