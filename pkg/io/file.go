package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WriteLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot create/truncate file with path '%s'", filePath)
	}
	defer file.Close()

	fileContent := strings.Join(lines, "\n")
	_, err = file.WriteString(fileContent)
	if err != nil {
		return fmt.Errorf("cannot write content to file with path '%s'", filePath)
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("cannot write content of file with path '%s' to disk", filePath)
	}

	return nil
}

func ReadLineByLine(filePath string, skipEmpty bool) (*chan string, error) {
	lines := make(chan string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file with path '%s'", filePath)
	}

	go func() {
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			line := fileScanner.Text()
			if skipEmpty && strings.TrimSpace(line) == "" {
				continue
			}
			lines <- line
		}
		close(lines)
		file.Close()
	}()

	return &lines, nil
}

func ReadBlocksOfLines(filePath string, skipEmpty bool, blockSize int) (*chan []string, error) {
	lines, err := ReadLineByLine(filePath, skipEmpty)
	if err != nil {
		return nil, err
	}
	blocks := make(chan []string)
	createBlock := func() []string {
		return make([]string, 0, blockSize)
	}

	go func() {
		block := createBlock()
		for line := range *lines {
			block = append(block, line)
			if len(block) == blockSize {
				blocks <- block
				block = createBlock()
			}
		}
		if len(block) > 0 {
			blocks <- block
		}
		close(blocks)
	}()

	return &blocks, nil
}
