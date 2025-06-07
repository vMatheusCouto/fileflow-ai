package fileutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ListFiles() ([]string, error) {
	var array = []string{}

	visit := func(path string, di fs.DirEntry, err error) error {
		newPath := strings.Replace(path, "../", "", -1)
		array = append(array, newPath)
		return nil
	}

	filepath.WalkDir("./files/", visit)
	files := array

	return files, nil
}

func MoveFile(from string, to string) {
	os.Rename("./files/"+from, "./result/"+to)
}
