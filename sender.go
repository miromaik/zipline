package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const chunkSize = 64 * 1024

func sendFile(path string) error {
	code := generateCode()
	key := deriveKey(code)

	fmt.Printf("\nTransfer Code: %s\n\n", code)
	fmt.Println("Waiting for receiver to connect...")

	conn, err := connectToRelay(code)
	if err != nil {
		return err
	}
	defer conn.Close()

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	var tmpFile string
	var size int64
	var filename string

	if info.IsDir() {
		tmpFile, err = zipDirectory(path)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile)

		zipInfo, _ := os.Stat(tmpFile)
		size = zipInfo.Size()
		filename = filepath.Base(path) + ".zip"
	} else {
		tmpFile = path
		size = info.Size()
		filename = filepath.Base(path)
	}

	msg := Message{
		Type: "init",
		Info: fmt.Sprintf("%s|%d", filename, size),
	}
	if err := sendMessage(conn, msg); err != nil {
		return err
	}

	response, err := receiveMessage(conn)
	if err != nil {
		return err
	}

	if !response.Confirm {
		fmt.Println("\nTransfer rejected by receiver")
		return nil
	}

	fmt.Println("Transfer accepted. Sending...")

	file, err := os.Open(tmpFile)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := newProgressBar(size)
	buf := make([]byte, chunkSize)

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		encrypted, err := encrypt(buf[:n], key)
		if err != nil {
			return err
		}

		msg := Message{
			Type: "data",
			Data: encrypted,
		}
		if err := sendMessage(conn, msg); err != nil {
			return err
		}

		bar.add(int64(n))
	}

	sendMessage(conn, Message{Type: "done", Done: true})
	bar.finish()

	fmt.Println("\nTransfer complete!")
	return nil
}

func zipDirectory(dirPath string) (string, error) {
	tmpFile, err := os.CreateTemp("", "zipline-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	baseDir := filepath.Base(dirPath)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(baseDir, relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}

		return nil
	})

	return tmpFile.Name(), err
}
