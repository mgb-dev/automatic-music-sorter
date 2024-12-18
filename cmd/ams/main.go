package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mgb-dev/ams/internal/metadata"
	"github.com/mgb-dev/ams/internal/utils"
)

func main() {
	// TODO: setup a proper argument parser
	if len(os.Args) < 3 {
		fmt.Println("Missing argument: pos 1 - 'work_dir'")
		fmt.Println("Missing argument: pos 2 - 'criteria'")
		return
	}
	workingDir := os.Args[1]
	criteria := os.Args[2]

	// in this case "debug means: Run simulated and verbosely"
	isDebugActive := len(os.Args) >= 4 && os.Args[3] == "--debug"
	if !metadata.IsValidCriteria(criteria) {
		log.Fatal("Invalid tag: ", criteria)
	}
	dirEntry, err := os.ReadDir(workingDir)
	if err != nil {
		log.Fatal("Directory reading error: \n", err)
	}
	fmt.Println("Working Directory :", workingDir)

	failures := 0
	successes := 0
	fileCount := len(dirEntry)

	dirList := utils.NewDirList()

	for _, fileEntry := range dirEntry {
		if fileEntry.IsDir() {
			continue
		}

		filename := fileEntry.Name()

		fmt.Println("File: ", filename)
		filePath := path.Join(workingDir, fileEntry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			failures++
			fmt.Println("File opening error:\n", err)
			continue
		}

		m, err := metadata.ReadTags(file)
		if err != nil {
			failures++
			fmt.Println("Metadata parsing error: ", err)
			continue
		}
		defer file.Close()

		tagData, ok := (*m.Raw())[criteria]
		if !ok || tagData == "" {
			failures++
			fmt.Printf("criteria: %s isn't available. Skipping file %s\n", criteria, filename)
			continue
		}
		newDirectory := path.Join(workingDir, utils.NormalizeDirName(tagData))
		newFilePath := path.Join(newDirectory, fileEntry.Name())

		if !dirList.Exists(newDirectory) {
			if isDebugActive {
				fmt.Printf("Adding to  DirList: % v\n", newDirectory)
			} else {
				if err := os.Mkdir(newDirectory, os.ModePerm); err != nil {
					// This error causes program crash as a security measure
					log.Fatal("Directory Creating error: ", err)
				}
			}
			dirList.Add(newDirectory)
		}

		if isDebugActive {
			fmt.Printf("moving %s -> %s\n", filePath, newFilePath)
		} else {
			if err := os.Rename(filePath, newFilePath); err != nil {
				failures++
				fmt.Println("File Relocation error: ", err)
				continue
			}
		}

		successes++
	}

	fmt.Printf(
		"Finished: %v file/s moved, %v file/s unchanged, Total files: %v\n",
		successes,
		failures,
		fileCount,
	)
}
