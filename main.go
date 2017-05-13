package main

import "os"
import "fmt"
import "bufio"
import "io/ioutil"
import "text/template"
import "encoding/json"

func main() {
	paths := os.Args[1:]

	if len(paths) != 2 {
		fmt.Println("Usage: consul-watch-renderer <template> <destination>")
		os.Exit(1)
	}

	sourcePath := paths[0]

	raw_template, error := ioutil.ReadFile(sourcePath)
	if error != nil {
		fmt.Println("Failed to read from template file:", sourcePath)
		os.Exit(2)
	}

	contents, error := template.New("template").Parse(string(raw_template))
	if error != nil {
		fmt.Println("Error parsing template:", error)
	}

	var data []interface{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		error := json.Unmarshal(scanner.Bytes(), &data)
		if error != nil {
			fmt.Println("Failed to parse stdin as json:", error)
			os.Exit(3)
		}
	}

	tempFile, error := ioutil.TempFile("/tmp", "cwr")
	if error != nil {
		fmt.Println("Failed to create temp file:", tempFile.Name(), error)
		os.Exit(4)
	}

	if error = contents.Execute(tempFile, data); error != nil {
		fmt.Println("Failed to render template:", error)
		os.Exit(5)
	}

	targetPath := paths[1]

	if error = os.Rename(tempFile.Name(), targetPath); error != nil {
		fmt.Println("Failed to replace target with rendered template:", tempFile.Name(), targetPath)
		os.Exit(6)
	} else {
		fmt.Printf("Successfully rendered template %s to %s\n", sourcePath, targetPath)
	}
}
