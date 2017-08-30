package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Return the AST as an array of strings, with the colorization removed
func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

// Start begins parsing an input file.
func Start(inputFile string, outputFile string, moduleName string) error {

	_, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("Input file is not found")
	}

	// Set up the info
	outputFilePath, outputFileName := filepath.Split(outputFile)
	modInfo := module{
		ModuleName:     moduleName,
		SourceHeader:   inputFile,
		OutputDir:      outputFilePath,
		OutputCFile:    outputFileName,
		OutputMakefile: "Makefile",
		OutputGoFile:   fmt.Sprintf("%sloader", moduleName),
	}

	// Preprocess
	var ppFilePath string

	var pp []byte
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	inputFilePath, _ := filepath.Split(inputFile)
	cmd := exec.Command("clang", "-E", "-I", inputFilePath, "-I", ".", inputFile)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
	}
	pp = []byte(out.String())

	tmpDir := os.TempDir()
	ppFilePath = path.Join(tmpDir, "pp.c")

	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return fmt.Errorf("writing to %s%cpp.c failed: %v", tmpDir, os.PathSeparator, err)
	}

	// Generate AST from preprocessed file
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", "-I", inputFilePath, "-I", ".", ppFilePath).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", ppFilePath).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	lines := readAST(astPP)

	modInfo.Functions = make([]parsedFunctionDefinition, 0)
	index := 0

	// Parse functions and their parameters
	inFunction := false
	functionRegex := regexp.MustCompile(`FunctionDecl.*col:\d+ (?P<funcname>\w+) '(?P<rettype>[\w\s*]+)\(`)
	functionAttributeRegex := regexp.MustCompile(`__attribute__\(\((?P<attribute>\w+)\)\)`)
	functionParamRegex := regexp.MustCompile(`ParmVarDecl.*col:\d+ (?P<name>\w+) '(?P<type>[\w\s*]+)'`)
	for _, line := range lines {
		if inFunction {
			paramInfo := functionParamRegex.FindStringSubmatch(line)
			if len(paramInfo) > 2 {
				modInfo.Functions[index].Parameters = append(modInfo.Functions[index].Parameters, parameterDefinition{name: paramInfo[1], dataType: paramInfo[2]})
			} else {
				inFunction = false
				index++
			}
		}
		if !inFunction {
			funcMatches := functionRegex.FindStringSubmatch(line)
			if len(funcMatches) > 2 {
				inFunction = true
				modInfo.Functions = append(modInfo.Functions, parsedFunctionDefinition{Name: funcMatches[1],
					ReturnType: strings.TrimSpace(funcMatches[2])})
				attribute := functionAttributeRegex.FindStringSubmatch(line)
				if len(attribute) > 1 {
					modInfo.Functions[index].Attribute = attribute[1]
				}
			}
		}
	}

	// Output new shim file
	err = modInfo.writeShimFile()
	if err != nil {
		return err
	}

	// Output new Makefile
	err = modInfo.writeMakefile()
	if err != nil {
		return err
	}

	// Output go stub files
	err = modInfo.writeGofiles()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s input_header_file output_c_file module_name\n", os.Args[0])
		os.Exit(1)
	}
	// Do the work
	if err := Start(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
