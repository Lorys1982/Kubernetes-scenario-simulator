package utils

import (
	"compress/gzip"
	"log"
	"main/configs"
	"os"
	"path"
	"strings"
	"sync"
)

var fileMutex sync.Mutex

// fileReplace Func to replace strings inside of files
func fileReplace(fileName string, toReplace string, replace string, input ...byte) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if len(input) == 0 {
		var err error
		input, err = os.ReadFile(fileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	output := strings.Replace(string(input), toReplace, replace, 1)

	err := os.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Compress func to compress files, taken in input fileName (name of the file to compress) and filePath (path to the dir
// of the file to compress)
//
// # Output
//
// The compressed file will be a .gz
func Compress(fileName string, filepath string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	workingDir, _ := os.Getwd()

	inFile := path.Join(filepath, fileName)
	input, _ := os.ReadFile(inFile)
	fileName = configs.GetCommandsName() + "_" + configs.LogTime + ".gz"

	outFile := path.Join(workingDir, "logs", fileName)
	newFile, _ := os.Create(outFile)
	compressor := gzip.NewWriter(newFile)
	_, err := compressor.Write(input)
	if err != nil {
		log.Println(err.Error())
	}
	err = compressor.Close()
	if err != nil {
		log.Println(err.Error())
	}
	defer newFile.Close()
}
