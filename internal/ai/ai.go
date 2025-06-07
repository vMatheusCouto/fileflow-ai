package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

type Folder struct {
	Name        string   `json:"name"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Files       []string `json:"files"`
}

type FoldersResponse struct {
	Folders map[string]Folder `json:"folders"`
}

type FilesResponse map[string][]string

type AIResponse struct {
	Folders *FoldersResponse
	Files   *FilesResponse
}

func CreateFolders(files []string) (string, error) {
	// Folder template
	folderConfig, errF := os.ReadFile("./config/folder.json")
	if errF != nil {
		log.Fatal(errF)
	}

	// API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	// Gemini initialization
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	prompt := fmt.Sprintf(`
	You are an expert in file management and organization.
	TASK: Create a folder structure for those files: %v.

	Template for structurizing: %s
	The maximum number of folders root folder can have is 8, so you may need to create subfolders (which you just add as another folder, only adapting the path). If there isnt too much files, you can create less than 8 root folders and also subfolders arent needed (but they can help a lot to organize).

	Analyze file content/context (path, name, extension, files nearby, etc).
	You should return ONLY a JSON object with the folder structure, based on ONLY the template.
	`, files, string(folderConfig))

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	resultStr := result.Text()

	return resultStr, err
}

func AssignFiles(files []string, folders string) (string, error) {
	// Files template
	filesConfig, err := os.ReadFile("./config/files.json")
	if err != nil {
		log.Fatal(err)
	}

	// Api key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	// Gemini initialization
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	prompt := fmt.Sprintf(`
	You are an expert in file management and organization.
	TASK: Assign those files to the right folders, based on the whole context (path, name, extension, files nearby, etc.): %v.

	Template for structurizing: %s
	This is the folder structure that you should assign to: %s

	IMPORTANT: the filename will be THE SAME that I sent to you, DO NOT CHANGE NOTHING.
	You should return ONLY a JSON object with the files structure, based on the template.
	`, files, string(filesConfig), folders)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	resultStr := result.Text()

	return resultStr, err
}
