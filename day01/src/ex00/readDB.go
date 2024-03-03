package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Recipes struct {
	Cakes []Cake `json:"cake" xml:"cake"`
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	StoveTime   string       `json:"time" xml:"stovetime"`
	Ingredients []Ingredient `json:"ingredients" xml:"ingredients>item"`
}

type Ingredient struct {
	Name  string `json:"ingredient_name" xml:"itemname"`
	Count string `json:"ingredient_count" xml:"itemcount"`
	Unit  string `json:"ingredient_unit" xml:"itemunit"`
}

type DBReader interface {
	Read(fileContent []byte) (Recipes, error)
}

type JSONReader struct{}

func (r JSONReader) Read(fileContent []byte) (Recipes, error) {
	var recipes Recipes
	err := json.Unmarshal(fileContent, &recipes)
	return recipes, err
}

type XMLReader struct{}

func (r XMLReader) Read(fileContent []byte) (Recipes, error) {
	var recipes Recipes
	err := xml.Unmarshal(fileContent, &recipes)
	return recipes, err
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbFilename := flag.String("f", "", "Read file by path, need to be .json or .xml")
	flag.Parse()

	if *dbFilename != "" {
		var reader DBReader
		var output []byte

		fileContent, err := os.ReadFile(*dbFilename)
		handleError(err)

		switch strings.ToLower(filepath.Ext(*dbFilename)) {
		case ".json":
			reader = JSONReader{}
		case ".xml":
			reader = XMLReader{}
		default:
			log.Fatal("Unsupported file extension")
		}

		recipes, err := reader.Read(fileContent)
		handleError(err)

		switch reader.(type) {
		case JSONReader:
			output, err = xml.MarshalIndent(recipes, "", "    ")
			handleError(err)
		case XMLReader:
			output, err = json.MarshalIndent(recipes, "", "    ")
			handleError(err)
		}
		fmt.Println(string(output))
	}
}
