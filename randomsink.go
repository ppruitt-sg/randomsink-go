package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var amount int
var fileType string
var fileName string

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&amount, "amount", 10, "The amount of random sink addresses")
	flag.StringVar(&fileType, "type", "json", "Type of file to be output (json or csv)")
	flag.StringVar(&fileName, "filename", "sink", "Name of csv or json file to be created)")
	flag.Parse()

	fileType = strings.ToLower(fileType)
}
func main() {
	var emails [][]byte
	err := validateFileType(fileType)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < amount; i++ {
		emails = append(emails, randomAddress())
	}

	if fileType == "json" {
		toJSON(emails, fileName)
	} else {
		toCSV(emails, fileName)
	}
}

func validateFileType(fileType string) error {
	if fileType != "csv" && fileType != "json" {
		return errors.New("Invalid file type")
	}
	return nil
}

func randomAddress() []byte {
	name := make([]byte, 10)
	for i := 0; i < 10; i++ {
		name[i] = charset[rand.Intn(25)]
	}
	return append(name, []byte("@sink.sendgrid.net")...)
}

func toCSV(emails [][]byte, fileName string) {
	f, err := os.Create(fileName + ".csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.Write([]byte("email\n"))

	for _, email := range emails {
		f.Write(append(email, []byte("\n")...))
	}

}

func toJSON(emails [][]byte, fileName string) {
	var emailsString []string
	f, err := os.Create(fileName + ".json")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	for i := 0; i < len(emails); i++ {
		emailsString = append(emailsString, string(emails[i]))
	}

	jsonMap := map[string]([]string){"email": emailsString}
	jsonString, err := json.MarshalIndent(jsonMap, "", "    ")

	if err != nil {
		log.Fatalln(err)
	}

	f.Write(jsonString)
}
