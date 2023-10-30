package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func fetchType(path string, ext string, grep string) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", ext, err)
		return
	}

	for _, entry := range dirs {
		if ext != "" {
			if entry.IsDir() || filepath.Ext(entry.Name()) != "."+ext {
				continue
			}
			filePath := filepath.Join(path, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s, %v\n", filePath, err)
			} else {
				// Check if the content contains the grep string
				if grep == "" || strings.Contains(string(content), grep) {
					fmt.Println(string(content))
					fmt.Println("#######################################")
				}
			}
		} else {
			if entry.IsDir() {
				continue
			}
			filename := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))]
			headerPath := filepath.Join(path, filename+".headers")
			bodyPath := filepath.Join(path, filename+".body")
			headerContent, err1 := os.ReadFile(headerPath)
			bodyContent, err2 := os.ReadFile(bodyPath)
			if err1 != nil {
				fmt.Printf("Error reading header file: %s\n", err1)
				return
			}

			if err2 != nil {
				fmt.Printf("Error reading body file: %s\n", err2)
				return
			}

			// Check if the header or body content contains the grep string
			if grep == "" || strings.Contains(string(headerContent), grep) || strings.Contains(string(bodyContent), grep) {
				fmt.Println(string(headerContent))
				fmt.Println()
				fmt.Println(string(bodyContent))
				fmt.Println("#######################################")
			}
			break
		}
	}
}

func main() {
	header := flag.Bool("h", false, "Fetch headers of fff saved outputs")
	body := flag.Bool("b", false, "Fetch request body of fff saved outputs")
	path := flag.String("p", ".", "Path where fff saved output is located. Default is current path")
	grep := flag.String("g", "", "Grep only the results containing the specified string")

	flag.Parse()
	fmt.Println("#######################################")

	// Loop through each item in the "path" directory and search for every directory inside it (not files)
	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			var ext string
			if *header && *body {
				ext = ""
			} else if *header {
				ext = "headers"
			} else if *body {
				ext = "body"
			} else {
				ext = ""
			}

			fetchType(path, ext, *grep)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
