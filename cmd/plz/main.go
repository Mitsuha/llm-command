package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s❌ Usage: plz <description>%s\n", ColorRed, ColorReset)
		fmt.Printf("%s💡 Example: plz show disk usage%s\n", ColorCyan, ColorReset)
		os.Exit(1)
	}

	// Start background cleanup
	go cleanupOldSessions()

	// Join all arguments after "plz" as the description
	description := strings.Join(os.Args[1:], " ")

	// Start loading animation
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)

	go showLoading(ctx, done)

	// Get command from OpenAI
	command, err := getCommandFromOpenAI(description)
	cancel() // Stop loading animation
	<-done   // Wait for loading animation to finish

	if err != nil {
		fmt.Printf("\r%s❌ Error: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	// Clear the thinking line and show command with prompt
	fmt.Printf("\r%s%s%s %s(Yes/No):%s ", ColorCyan, command, ColorReset, ColorYellow, ColorReset)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	accepted := response == "y" || response == "yes"

	// Save session entry
	saveSessionEntry(description, command, accepted)

	if accepted {
		executeCommand(command)
	}
}
