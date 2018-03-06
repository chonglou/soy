package main

import (
	"log"
	"os"

	_ "github.com/chonglou/soy/blog"
	_ "github.com/chonglou/soy/cbeta"
	_ "github.com/chonglou/soy/dict"
	"github.com/chonglou/soy/env"
)

func main() {
	if err := env.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
