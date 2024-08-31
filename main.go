package main

import (
	"os"

	"github.com/Terrorknubbel/gitmate/internal/cmd"
)

func main() {
	cmd.Main(os.Args[1:])
}
