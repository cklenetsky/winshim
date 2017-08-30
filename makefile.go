package main

import (
	"io"
	"strings"
	"text/template"
)

type makefileInfo struct {
	ModuleName        string
	FileNameBase      string
	ConvertedFilePath string
}

func initMakefileTemplate() (*template.Template, error) {
	funcMap := template.FuncMap{
		"tolower": strings.ToLower,
	}
	return template.New("makefile.template").Funcs(funcMap).ParseFiles("makefile.template")
}

func writeMakefile(moduleName string, inputFilePath string, outputFileName string, output io.Writer) error {

	makefileTemplate, err := initMakefileTemplate()
	if err != nil {
		return err
	}

	// Split outputfilename
	var outputFileNameBase string
	i := strings.LastIndex(outputFileName, ".")
	if i > 0 {
		outputFileNameBase = outputFileName[:i]
	} else {
		outputFileNameBase = outputFileName
	}

	// For MinGW, the filepath must be converted
	convertedFilePath := strings.Replace(inputFilePath, "\\", "/", -1)

	infoStruct := makefileInfo{
		ModuleName:        moduleName,
		FileNameBase:      outputFileNameBase,
		ConvertedFilePath: convertedFilePath,
	}
	err = makefileTemplate.ExecuteTemplate(output, "makefile.template", infoStruct)
	return err
}
