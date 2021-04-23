package main

import (
	"os"

	"github.com/shayan-7/gocha/internal/cmd"
)

func main() {
	args := os.Args
	cmd.Execute(args)
}
