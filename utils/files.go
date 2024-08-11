package utils

import (
	"log"
	"os"
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
