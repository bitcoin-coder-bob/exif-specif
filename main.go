package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	var lastViewerProcess *os.Process // Track the last viewer process

	// Open the directory and read files
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error reading directory: %s", err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		// Close the previous image viewer, if any
		if lastViewerProcess != nil {
			err := lastViewerProcess.Kill()
			if err != nil {
				log.Printf("Error closing previous image viewer: %s", err.Error())
			}
			lastViewerProcess = nil
		}

		// Open the image directly with eog
		cmd := exec.Command("eog", filePath)
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Error opening file %s with eog: %s", filePath, err.Error())
		}
		lastViewerProcess = cmd.Process // Track the eog process

		// Extract current keywords
		cmd = exec.Command("exiftool", "-Keywords", "-sep", ", ", filePath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Error reading keywords from file %s: %s", filePath, err.Error())
		}

		// Clean up the output to remove the "Keywords :" label
		rawOutput := out.String()
		rawOutput = strings.ReplaceAll(rawOutput, "Keywords                        :", "") // Remove label
		rawOutput = strings.TrimSpace(rawOutput)                                           // Remove leading/trailing spaces

		// Parse current keywords into a slice and clean up
		rawKeywords := strings.Split(rawOutput, ", ")
		currentKeywords := []string{}
		for _, keyword := range rawKeywords {
			trimmedKeyword := strings.TrimSpace(keyword)
			if trimmedKeyword != "" { // Exclude empty keywords
				currentKeywords = append(currentKeywords, trimmedKeyword)
			}
		}

		// Prompt user for new keywords
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter tags for %s (space-separated): ", file.Name())
		keywords, _ := reader.ReadString('\n')
		keywords = strings.TrimSpace(keywords)
		newKeywords := strings.Fields(keywords) // Split by spaces

		// Add only new keywords that aren't already present
		for _, newKeyword := range newKeywords {
			if !contains(currentKeywords, newKeyword) {
				currentKeywords = append(currentKeywords, newKeyword)
			}
		}

		// Construct the final keywords argument
		finalKeywords := strings.Join(currentKeywords, ", ")
		fmt.Printf("Final keywords for %s: %s\n", file.Name(), finalKeywords)

		// Add updated keywords using exiftool
		keywordArg := fmt.Sprintf("-Keywords=%s", finalKeywords)
		cmd = exec.Command("exiftool", "-overwrite_original", keywordArg, filePath)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error tagging file %s: %s\n", filePath, err.Error())
			log.Fatal(err)
		}
	}

	// Close the last image viewer, if any
	if lastViewerProcess != nil {
		err := lastViewerProcess.Kill()
		if err != nil {
			log.Printf("Error closing final image viewer: %s", err.Error())
		}
	}

	fmt.Println("All images have been tagged.")
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
