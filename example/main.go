package main

import (
	"fmt"
	"os"

	"github.com/fpabl0/zipper-go"
)

func main() {

	const tempPath = "temp/123"

	err := os.RemoveAll(tempPath)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(tempPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	rootFolders, filenames, err := zipper.Unzip("temp.zip", tempPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("rootFolders:", rootFolders)
	fmt.Println("Files:", len(filenames))
}
