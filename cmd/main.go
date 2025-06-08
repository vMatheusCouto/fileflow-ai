package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"fileflow-ai/internal/ai"
	"fileflow-ai/internal/fileutils"
	"fileflow-ai/internal/folderutils"

	"github.com/joho/godotenv"
	"github.com/yarlson/pin"
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
	fmt.Print("→ Files organization started! \n\n")
	godotenv.Load(".env.local")
	scanner := bufio.NewScanner(os.Stdin)

	// Load files
	files, err := fileutils.ListFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Insert the language to be utilized → " + "\033[31m")
	scanner.Scan()
	language := scanner.Text()

	fmt.Print("\033[30m" + "Extra instructions → " + "\033[31m")
	scanner.Scan()
	category := scanner.Text()

	folderS := pin.New("Creating folder structure...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
	)
	cancelFolderS := folderS.Start(context.Background())
	defer cancelFolderS()

	folderStructure, err := ai.CreateFolders(files, language, category)
	if err != nil {
		fmt.Println(err)
	}

	converted := strings.Replace(strings.Replace(folderStructure, "```json", "", -1), "```", "", -1)

	var foldersResp FoldersResponse
	errJ := json.Unmarshal([]byte(converted), &foldersResp)
	if errJ != nil {
		fmt.Println(err)
	}

	for _, folder := range foldersResp.Folders {
		folderutils.CreateFolder(folder.Path)
	}
	folderS.Stop("Created folders sucessfully!")
	fmt.Print("\n")

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
		message := "Assigning files " + strconv.Itoa(j+1) + "/" + strconv.Itoa(len(newResponse)) + "..."
		filesA := pin.New(message,
			pin.WithSpinnerColor(pin.ColorCyan),
			pin.WithTextColor(pin.ColorYellow),
		)
		cancelFilesA := filesA.Start(context.Background())
		defer cancelFilesA()
		result, err := ai.AssignFiles(newResponse[j], folderStructure)
		if err != nil {
			fmt.Println(err)
		}
		converted2 := strings.Replace(strings.Replace(result, "```json", "", -1), "```", "", -1)

		// Run files moving
		var filesResp FilesResponse
		err2 := json.Unmarshal([]byte(converted2), &filesResp)
		if err2 != nil {
			fmt.Println(err2)
		}
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
		filesA.Stop("Assigned files " + strconv.Itoa(j+1) + "/" + strconv.Itoa(len(newResponse)) + ".")
	}

	os.Rename("./files/", "./trash/")
	os.Mkdir("./files/", os.ModePerm)
	fmt.Print("\n")
	fmt.Println("\033[31m" + "→" + "\033[0m" + "\033[34m" + " Sucessfully completed organization!")
	fmt.Println("\033[34m" + "Files: " + "\033[0m" + "result/")
	fmt.Println("\033[34m" + "Remaining " + "\033[0m" + "trash/")
}
