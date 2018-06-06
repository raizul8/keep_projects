package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func renameFiles(root, oldPath, newName string) error {
	newPath := fmt.Sprintf("%s/%s", root, newName)

	err := exec.Command("mv", oldPath, newPath).Run()
	if err != nil {
		return err
	}

	return nil
}

func readRenameFiles(root string, fileNameStart string) {
	nr := 0
	if fileNameStart == "" {
		return
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		newName := fmt.Sprintf("%s (%s)%s", fileNameStart, strconv.Itoa(nr), ".jpg")
		nr++

		err = renameFiles(root, path, newName)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("renameFiles-> %v", err)
		return
	}

}

func main() {
	fmt.Println("--------Start --------------")

	root := "/home/raz/Pictures/bl1"
	//root := "/home/raz/IdeaProjects1/keep_projects/rename_files/src/div"
	readRenameFiles(root, "bla")

}
