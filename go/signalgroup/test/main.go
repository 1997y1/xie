package main

import (
	"fmt"
	"time"

	"xie/go/signalgroup"
)

func main() {
	signalgroup.Async(0, func() error {
	loop:
		time.Sleep(100 * time.Millisecond)
		fmt.Println(time.Now())
		goto loop
	})
	signalgroup.Async(0, func() error {
		time.Sleep(5 * time.Second)
		return nil
	})
	signalgroup.Wait()

	fmt.Println("Ending")
}
