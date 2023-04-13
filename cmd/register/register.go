package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gofmt "go/format"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/actgardner/gogen-avro/v10/generator"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("too few arguments provided: [path] [avro-filename] [template-id]")
	}

	path := os.Args[1]
	avroFilename := os.Args[2]
	templateID := os.Args[3]

	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	avroPath := filepath.Join(path, avroFilename)
	avroSchema := loadAvroSchema(avroPath)

	editGeneratedFile(path, templateID, avroSchema.Name)
}

type AvroSchema struct {
	Name string
}

func loadAvroSchema(path string) AvroSchema {
	avroFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var avroSchema AvroSchema
	err = json.NewDecoder(avroFile).Decode(&avroSchema)
	if err != nil {
		log.Fatal(err)
	}

	return avroSchema
}

func editGeneratedFile(path, templateID, typeName string) {
	initFilename := filepath.Join(path, generator.ToSnake(typeName)+".go")

	// final result
	outBuffer := &bytes.Buffer{}

	// save the top comment (AST will drop those)
	content, err := os.ReadFile(initFilename)
	if err != nil {
		log.Fatal(err)
	}
	comments := strings.Split(string(content), "package")[0]
	outBuffer.WriteString(comments)

	// open the file and handle the AST
	fset := token.NewFileSet()
	initFile, err := parser.ParseFile(fset, initFilename, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	// add the import
	astutil.AddImport(fset, initFile, `github.com/ugol/jr/types/registry`)
	printer.Fprint(outBuffer, fset, initFile)

	// add the init func at the bottom
	template := `
func init() {
	registry.Register("%s", &%s{})
}
`
	_, err = fmt.Fprintf(outBuffer, template, templateID, typeName)
	if err != nil {
		log.Fatal(err)
	}

	// go fmt the file
	goFormattedContent, err := gofmt.Source(outBuffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// write out the result
	outFile, err := os.OpenFile(initFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	_, err = outFile.Write(goFormattedContent)
	if err != nil {
		log.Fatal(err)
	}
	outFile.Close()
}
