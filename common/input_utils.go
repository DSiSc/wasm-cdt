package common

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"runtime"
	"syscall"
)

// GetPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func GetPassPhrase(prompt string, confirmation bool) string {
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		Fatalf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		fmt.Print("\nRepeat Password: ")
		bytePasswordConfirm, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if string(bytePassword) != string(bytePasswordConfirm) {
			Fatalf("Passphrases do not match")
			return ""
		}
		fmt.Println("")
	}
	return string(bytePassword)
}

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}
