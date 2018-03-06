package main

import (
	"log"
	"os"

	"github.com/chonglou/soy/env"
)

func main() {
	if err := env.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
