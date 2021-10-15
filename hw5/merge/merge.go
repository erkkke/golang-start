package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func generator(val int) <-chan int {
	ch := make(chan int)

	go func() {
		for i := 1; i <= val; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	channel := make(chan int)
	

	for _, c := range cs {
		wg.Add(1)
		
		go func(ch <-chan int) {
			defer wg.Done()
			for in := range ch {
				channel <- in
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	return channel
}

func main() {
	var cs []<-chan int
	for i := 0; i < 6; i++ {
		cs = append(cs, generator(rand.Intn(10)))
	}

	mergedChannel := merge(cs...)

	for val := range mergedChannel {
		fmt.Println(val)
	}

}