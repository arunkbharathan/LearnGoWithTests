package main

import (
	"os"
	"time"

	clockface "github.com/arunkbharathan/learnWithTests/maths"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
