package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func main1() {
	initArr := []string{"I", "am", "stupid", "and", "weak"}

	for i, v := range initArr {
		if v == "stupid" {
			initArr[i] = "smart"
		}

		if v == "weak" {
			initArr[i] = "strong"
		}
	}

	fmt.Println(strings.Join(initArr, " "))
}

func main() {
	que := make(chan int, 10)
	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	// producer
	go func() {
		for {
			select {
			case <-ticker.C:
				que <- 1
			case <-ctx.Done():
				fmt.Println("producer exit")
				return
			}
		}

	}()

	// consumer
	go func() {
		for {
			select {
			case v := <-que:
				fmt.Printf("consume data %d\n", v)
			case <-ctx.Done():
				fmt.Println("consumer exit")
				return
			}
		}
	}()

	time.Sleep(time.Second * 10)
	cancel()
	time.Sleep(time.Second * 2)

}
