package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func AssignFiles(files []string) ([]string, error) {
	folderConfig, errFo := os.ReadFile("./config/folder.json")
	filesConfig, errFi := os.ReadFile("./config/files.json")
	if errFo != nil || errFi != nil {
		log.Fatal("Error reading configuration files")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text("You are expert in file management. Create a folder structure for a project that includes the following files: "+fmt.Sprintf("%v", files)+". Return ONLY 2 json: the folder structure, based on "+string(folderConfig)+", and the json of the files that should be in each folder, based on"+string(filesConfig)+". You should analyse the whole context, and not only the extension of the file. Do not return any other text or explanation."),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	resultStr := result.Text()
	return []string{resultStr}, err
}
