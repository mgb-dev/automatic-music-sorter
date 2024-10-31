package main

import (
	"fmt"
	"log"
	"os"

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

	for _, fileEntry := range dirEntry {
		if fileEntry.IsDir() {
			continue
		}

		filename := fileEntry.Name()

		fmt.Println("File: ", filename)
		filePath := workingDir + fileEntry.Name()
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal("File opening error:\n", err)
		}

		m, err := metadata.ReadTags(file)
		if err != nil {
			log.Fatal("Metadata parsing error: ", err)
		}
		defer file.Close()

		tagData, ok := (*m.Raw())[criteria]
		if !ok {
			fmt.Printf("criteria: %s isn't available. Skipping file %s", criteria, filename)
			continue
		}
		newDirectory := workingDir + tagData
		fmt.Printf("New dir: %s => %s\n", newDirectory, filename)

		// if err := os.Mkdir(newDirectory, os.ModeDir); err != nil {
		// 	log.Fatal("Directory Creating error: ", err)
		// }

		// if err := os.Rename(filePath, (newDirectory + fileEntry.Name())); err != nil {
		// 	log.Fatal("File Relocation error: ", err)
		// }

	}

	fmt.Println("Program terminated successfully")
}
