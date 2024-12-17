package utils

import (
	"compress/gzip"
	"log"
	"main/global"
	"os"
	"path"
	"strings"
	"sync"
)

var fileMutex sync.Mutex

// FileReplace Func to replace strings inside of files
//
// # Parameters
//
// filename: name of the file in which replace the string
//
// toReplace: the string to change
//
// replace: the replacement string
//
// input: if you have already opened and read the file
// you can send it directly (if you have to do many iterations this is faster)
func FileReplace(fileName string, toReplace string, replace string, input Option[[]byte]) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if input.IsNone() {
		var err error
		res, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
		input = Some(res)
	}

	output := strings.Replace(string(input.GetSome()), toReplace, replace, 1)

	err := os.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Compress func to compress files, taken in input fileName (name of the file to compress) and filePath (path to the dir
// of the file to compress)
//
// # Result
//
// The compressed file will be a .gz
func Compress(fileName string, filepath string, clusterIndex int) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	workingDir, _ := os.Getwd()

	inFile := path.Join(filepath, fileName)
	input, _ := os.ReadFile(inFile)
	fileName = global.ConfName[clusterIndex] + "_" + global.LogTime + ".gz"

	outFile := path.Join(workingDir, "logs", global.ClusterNames[clusterIndex], fileName)
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

func CleanLogs() {
	os.RemoveAll("./logs/")
	os.MkdirAll("./logs", os.ModePerm)
}
