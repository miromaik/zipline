package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func receiveFile(code string) error {
	key := deriveKey(code)

	fmt.Println("Connecting to sender...")

	conn, err := connectToRelay(code)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg, err := receiveMessage(conn)
	if err != nil {
		return err
	}

	if msg.Type != "init" {
		return fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	var filename string
	var size int64
	fmt.Sscanf(msg.Info, "%s|%d", &filename, &size)

	location := getApproximateLocation()
	fmt.Printf("\nIncoming transfer request:\n")
	fmt.Printf("  File: %s\n", filename)
	fmt.Printf("  Size: %s\n", formatBytes(size))
	fmt.Printf("  From: %s\n\n", location)
	fmt.Print("Accept transfer? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	confirm := response == "yes" || response == "y"

	if err := sendMessage(conn, Message{Type: "confirm", Confirm: confirm}); err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Transfer rejected")
		return nil
	}

	fmt.Println("\nReceiving file...")

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bar := newProgressBar(size)

	for {
		msg, err := receiveMessage(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if msg.Done {
			break
		}

		if msg.Type == "data" {
			decrypted, err := decrypt(msg.Data, key)
			if err != nil {
				return fmt.Errorf("decryption failed: %w", err)
			}

			if _, err := outFile.Write(decrypted); err != nil {
				return err
			}

			bar.add(int64(len(decrypted)))
		}
	}

	bar.finish()
	fmt.Printf("\nFile saved: %s\n", filename)
	return nil
}

func getApproximateLocation() string {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return "Unknown Location"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown Location"
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "Unknown Location"
	}

	city, _ := result["city"].(string)
	country, _ := result["country_name"].(string)

	if city != "" && country != "" {
		return fmt.Sprintf("%s, %s", city, country)
	}
	if country != "" {
		return country
	}
	return "Unknown Location"
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
