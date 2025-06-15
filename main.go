package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Check for root privileges
	if os.Getuid() != 0 {
		log.Fatal("This program requires root privileges to run")
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: sudo ./phoenix-recovery <device_path>")
		fmt.Println("Example: sudo ./phoenix-recovery /dev/sda1")
		os.Exit(1)
	}

	devicePath := os.Args[1]
	fmt.Printf("Analyzing device: %s\n", devicePath)

	// Create filesystem analyzer
	analyzer, err := NewFilesystemAnalyzer(devicePath)
	if err != nil {
		log.Fatalf("Failed to create filesystem analyzer: %v", err)
	}
	defer analyzer.Close()

	// Read superblock information
	if err := analyzer.ReadSuperblock(); err != nil {
		log.Fatalf("Failed to read superblock: %v", err)
	}

	fmt.Println("Filesystem information:")
	analyzer.PrintFilesystemInfo()
}