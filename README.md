# winshim
A tool to generate a shim libarary for use with DLLs on Windows

## Introduction
When developing with cgo on Unix there is a fairly easy integration between go and existing libraries.  Using .so or .a files is relatively simple, with LDFLAGS able to use -l<library> to link appropriately.  On Windows, however, it may not be as simple.  Go does not natively link with Visual Studio-created .lib files, nor can it directly use .dlls.  It is possible to work around these limitations, and this tool serves to make that task easer.

winshim creates a series of files, for use on both Windows and non-Windows platforms, that allows for a cross-platform solution.  It generates a .c file and Makefile to build a small shim library that will dynamically load the target DLL, retrieve function pointers to the functions, and call the functions through those pointers.  It also generates to .go files, one for Windows, one for other platforms, that calls into the shim library, on Windows, or does nothing, on other platforms.

## Requirements
winshim relies on LLVM's clang to be installed and in the path.  clang is used to preprocess and parse the header file associated with the target library and to generate an AST used by winshim to create its files.  It also assumes, for the resulting shim library, that MinGW is installed, as this is the supported C integration system on Windows.

## Usage
d:>winshim.exe input_header_file output_c_file module_name

Parameter | Description
--------- | -----------
input_header_file | the main C header file containing the API for the DLL
output_c_file | the full path to the .c file to write.  All of the generated files will be written to this directory.
module_name | a name for this DLL/module.  Only alphanumeric characters are allowed, as the name is used in some of the generated functions.

## Results
File | Description
--------- | -----------
output_c_file | the MinGW-compatible shim library file.
Makefile | the Makefile used to build the shim library.
module_nameloader_windows.go | a .go file to be used with the integrating application.  Using cgo it calls into the shim library to initialize and clean up the shim library
module_nameloader.go | a .go file to be used with the integrating application.  On non-Windows platforms it provides no-op functions mirroring those provided by module_nameloader_windows.go

