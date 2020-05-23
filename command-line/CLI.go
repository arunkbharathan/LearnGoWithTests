package poker

import (
	"bufio"
	"io"
)

// CLI helps players through a game of poker
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

// NewCLI creates a CLI for playing poker
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}
