package main

import (
	"fmt"
	"os"
)

const version = "1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "send":
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}

		if err := sendFile(path); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "get", "receive":
		if len(os.Args) < 3 {
			fmt.Println("Error: code required")
			fmt.Println("Usage: zl get <code>")
			os.Exit(1)
		}

		code := os.Args[2]
		if len(code) != 6 {
			fmt.Println("Error: code must be 6 digits")
			os.Exit(1)
		}

		if err := receiveFile(code); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "relay":
		addr := ":8888"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}

		server := NewRelayServer()
		if err := server.Start(addr); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "version", "-v", "--version":
		fmt.Printf("Zipline v%s\n", version)

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`Zipline v%s - Secure P2P file transfer

Usage:
  zl send [file/folder]    Send file or folder
  zl get <code>            Receive file with 6-digit code
  zl relay [addr]          Run relay server
  zl version               Show version

Examples:
  zl send file.pdf
  zl send ./folder
  zl get 123456

`, version)
}
