package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const makefileContents = `# This Makefile assumes that MinGW is in use and in the path
WINLIB := lib%sshim.a
WINLIBDEP := %s.c
WINLIBOBJ := %s.o

CC := gcc
CCFLAGS := -std=c99

AR := ar
ARFLAGS := crf

$(WINLIB): $(WINLIBDEP)
	$(CC) $(CCFLAGS) -c $(WINLIBDEP) -o $(WINLIBOBJ) -I %s
	$(AR) $(ARFLAGS) $(WINLIB) $(WINLIBOBJ)
`

func writeMakefile(inputFilePath string, outputFileName string, output io.Writer) error {
	var b bytes.Buffer

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

	//
	b.WriteString(fmt.Sprintf(makefileContents, outputFileNameBase, outputFileNameBase, outputFileNameBase, convertedFilePath))

	_, err := output.Write(b.Bytes())
	return err
}
