package main

import (
	"os"

	"github.com/alecthomas/repr"
	"github.com/vinyl-linux/mint/parser"
)

func main() {
	dir := os.Getenv("DIR")

	a, err := parser.ParseDir(dir)
	if err != nil {
		panic(err)
	}

	repr.Println(a)
}
