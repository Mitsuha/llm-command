package main

import (
	"context"
	"fmt"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func showLoading(ctx context.Context, done chan bool) {
	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-ctx.Done():
			done <- true
			return
		default:
			fmt.Printf("\r%s🤔 Thinking... %s%s", ColorBlue, chars[i%len(chars)], ColorReset)
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}