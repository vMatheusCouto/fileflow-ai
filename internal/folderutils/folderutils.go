package folderutils

import (
	"fmt"
	"os"
)

func CreateFolder(path string) error {
	err := os.MkdirAll(string("./files"+path), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
