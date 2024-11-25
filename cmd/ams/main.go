package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mgb-dev/ams/internal/metadata"
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
		fmt.Printf("moving %s -> %s\n", filePath, newFilePath)

		// if err := os.Mkdir(newDirectory, os.ModeDir); err != nil {
		// 	log.Fatal("Directory Creating error: ", err)
		// }

		// if err := os.Rename(filePath, (newDirectory + fileEntry.Name())); err != nil {
		// 	log.Fatal("File Relocation error: ", err)
		// }

		successes++
	}

	fmt.Printf(
		"Finished: Successes: %v file/s moved. Failures: %v file/s unchanged. Total files: %v\n",
		successes,
		failures,
		fileCount,
	)
}
