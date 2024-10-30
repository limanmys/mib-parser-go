package mibparser

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Load initializes the MIBParser with options
func Load(opts ...Option) (*MIBParser, error) {
	opt := Opts{}
	for _, o := range opts {
		o(&opt)
	}
	return &MIBParser{opts: opt}, nil
}

// ReadMIBFile reads all MIB files in the directory
func (p *MIBParser) ReadMIBFile() ([]string, error) {
	files, err := os.ReadDir(p.opts.Path)
	if err != nil {
		log.Fatal(err)
	}
	var mergedMIB []string
	for _, file := range files {
		lines, err := readMIBFileWithPath(p.opts.Path + "/" + file.Name())
		if err != nil {
			fmt.Println("Error reading MIB file:", err) // Print error if a file can't be read
		}
		mergedMIB = append(mergedMIB, lines...) // Merge lines into a single slice
	}
	return mergedMIB, nil
}

// readMIBFileWithPath reads a single MIB file and returns its lines
func readMIBFileWithPath(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when function exits
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // Read each line into the slice
	}
	return lines, scanner.Err()
}
