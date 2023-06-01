package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func updateFileLine(filename string, currentNode int, nodeNum int) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Open file err:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Close file err:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var fileData string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "laddr = \"tcp://127.0.0.1:26657\"") {
			line = strings.Replace(line, "laddr = \"tcp://127.0.0.1:26657\"", "laddr = \"tcp://0.0.0.0:26657\"", -1)
		}
		if strings.Contains(line, "create-empty-blocks = true") {
			line = strings.Replace(line, "create-empty-blocks = true", "create-empty-blocks = false", -1)
		}

		if strings.Contains(line, "indexer = [\"null\"]") {
			line = strings.Replace(line, "indexer = [\"null\"]", "indexer = [\"kv\"]", -1)
		}

		if strings.Contains(line, "queue-type = ") {
			line = strings.Replace(line, "simple-priority", "priority", -1)
		}

		if strings.Contains(line, "max-subscription-clients ") {
			line = strings.Replace(line, "100", "1000000", -1)
		}

		if strings.Contains(line, "max-subscriptions-per-client =") {
			line = strings.Replace(line, "5", "1000000", -1)
		}

		if strings.Contains(line, "size = 5000") {
			line = strings.Replace(line, "size = 5000", "size = 50000", -1)
		}

		if strings.Contains(line, "cache-size = 10000") {
			line = strings.Replace(line, "cache-size = 10000", "cache-size = 100000", -1)
		}

		fileData += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	output, err := os.Create(filename)
	if err != nil {
		fmt.Println("Create err:", err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			fmt.Println("output Close err:", err)
		}
	}(output)

	_, err = output.WriteString(fileData)
	return err
}
