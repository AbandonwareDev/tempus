package main

import (
	"fmt"
	"os"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	options, err := ParseOptions()
	errHandler(err, "Error parsing options")

	fmt.Println("hello world", options)

	// start workers in parallel
	for i := 0; i < options.Threads; i++ {
		waitGroup.Add(1)
		go func() {
			fmt.Println("do parallel stuff")

			defer waitGroup.Done()
		}()
	}
	waitGroup.Wait()

}

func errHandler(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}
