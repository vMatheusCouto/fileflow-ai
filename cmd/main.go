package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fileflow-ai/internal/ai"
	"fileflow-ai/internal/fileutils"
	"fileflow-ai/internal/folderutils"

	"github.com/joho/godotenv"
)

type Folder struct {
	Name        string   `json:"name"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Files       []string `json:"files"`
}

type FoldersResponse struct {
	Folders map[string]Folder
}

type FilesResponse map[string][]string

func main() {
	fmt.Println("Files organization started! \n")
	godotenv.Load(".env.local")

	// Load files
	files, err := fileutils.ListFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create folder structure
	fmt.Println("Using AI to create folders structure...")
	folderStructure, err := ai.CreateFolders(files)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Folders structure done.")

	converted := strings.Replace(strings.Replace(folderStructure, "```json", "", -1), "```", "", -1)

	var foldersResp FoldersResponse
	errJ := json.Unmarshal([]byte(converted), &foldersResp)
	if errJ != nil {
		fmt.Println(err)
	}

	for _, folder := range foldersResp.Folders {
		folderutils.CreateFolder(folder.Path)
	}
	fmt.Println("Directories created sucessfully. \n")

	// Group files (50 each)
	newResponse := [][]string{}
	for i := 0; i < len(files); i += 50 {
		end := i + 50
		if end > len(files) {
			end = len(files)
		}
		newResponse = append(newResponse, files[i:end])
	}

	// Process each group of files
	for j := range newResponse {

		// Assign files
		fmt.Println("Using AI to assign file group " + string(j+1) + "...")
		result, err := ai.AssignFiles(newResponse[j], folderStructure)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Files assigned sucessfully.")
		converted2 := strings.Replace(strings.Replace(result, "```json", "", -1), "```", "", -1)

		// Run files moving
		var filesResp FilesResponse
		err2 := json.Unmarshal([]byte(converted2), &filesResp)
		if err2 != nil {
			fmt.Println(err2)
		}
		fmt.Println("Moving files to correct folder...")
		for folderPath, files := range filesResp {
			for _, file := range files {
				filePath, has := strings.CutPrefix(file, "files/")
				if has {
				}
				length := strings.Split(filePath, "/")
				name := length[len(length)-1]
				pathTo := foldersResp.Folders[string(folderPath)].Path + "/" + name

				fileutils.MoveFile(filePath, pathTo)
			}
		}
		fmt.Println("Files moved sucessfully. \n")
		os.Rename("./files/", "./trash/")
		os.Mkdir("./files/", os.ModePerm)
	}

	fmt.Println("Organization completed.")
	fmt.Println("Files are at -/result/.")
	fmt.Println("Remaining folders or files are at -/trash/.")
}
