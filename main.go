package main

import (
	"bufio"
	"fmt"
	"os"
)

type PageData struct {
	ExtenalLinks  []string
	InternalLinks []string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error.
	}
}
