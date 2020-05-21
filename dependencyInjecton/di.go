package main

import (
	"io"
	"os"
)

func Greet(writable io.Writer, val string) {
	writable.Write([]byte("Hello, " + val))
}

func main() {
	Greet(os.Stdout, "Elodie")
}
