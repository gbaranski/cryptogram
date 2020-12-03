package repl

import (
	"fmt"
	"os"
	"strings"
)

// Executor used to execute commands in CLI
func Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Goodbye")
		os.Exit(0)
		return
	}

	fmt.Println("Received " + s)

}
