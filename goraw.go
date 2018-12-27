package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"./libraw"
)

func main() {

	path := os.Args[1]
	exportPath := path + "/_export"
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		os.Mkdir(exportPath, 0777)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".ORF" {
			libraw.Export(path, f, exportPath)
		}

	}

}
