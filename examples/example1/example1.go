package main

import (
	"fmt"

	"github.com/romnnn/mongotypes"
)

func run() string {
	return mongotypes.Shout("This is an example")
}

func main() {
	fmt.Println(run())
}
