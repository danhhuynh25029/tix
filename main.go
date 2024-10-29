package main

import (
	"fmt"
	"log"
	"os"
	"tix/pkg"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run zip.go archive.zip file1 [file2...]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "zip":
		err := pkg.ZipFiles(os.Args[2], os.Args[3:]...)
		if err != nil {
			log.Fatal(err)
		}
	case "unzip":
		err := pkg.Unzip(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	}
}
