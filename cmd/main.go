package main

import (
	"fmt"
	"log"

	"fileflow-ai/internal/ai"
	"fileflow-ai/internal/fileutils"

	"github.com/joho/godotenv"
)

type Types struct {
}

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env:", err)
	}

	fmt.Println("Hello, World!")
	files, err := fileutils.ListFiles("./files")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Printf("Total files found: %d\n", len(files))
	fmt.Println("Files:", files)
	result, err := ai.AssignFiles(files)
	fmt.Println("AI Result:", result)
}
