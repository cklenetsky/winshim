package main

type module struct {
	ModuleName     string                     // The base name of the output module
	SourceHeader   string                     // The C header file used to create the shim
	Functions      []parsedFunctionDefinition // The functions parsed from the C header file
	OutputDir      string                     // The directory to write the generated files
	OutputCFile    string                     // The name of the .c file to generate
	OutputGoFile   string                     // The base name of the .go files to generate
	OutputMakefile string                     // The name of the Makefile to generate
}
