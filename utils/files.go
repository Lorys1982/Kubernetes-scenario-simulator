package utils

import (
	"compress/gzip"
	"log"
	"os"
	"path"
	"strings"
)

func fileReplace(fileName string, toReplace string, replace string, input ...byte) {
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

func Compress(fileName string, filepath string) {
	workingDir, _ := os.Getwd()

	inFile := path.Join(filepath, fileName)
	input, _ := os.ReadFile(inFile)
	fileName = strings.Replace(fileName, ".log", ".gz", 1)

	outFile := path.Join(workingDir, "Logs", fileName)
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
