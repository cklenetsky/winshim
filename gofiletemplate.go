package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const windowsGoStub = `// +build windows

package <INSERT PACKAGE NAME HERE>

// #cgo CPPFLAGS: -I <INSERT PATH TO INCLUDE FILE HERE>
// #cgo LDFLAGS: -L<INSERT PATH TO SHIM LIBRARY HERE> -lstdc++ -l%sshim
//
// Init%s(const char *libname);
// Cleanup%s(void);
import "C"

func loaderInit() {
	C.Init%s(C.CString("<INSERT DLL FILENAME HERE>"))
}

func loaderDestroy() {
	C.Cleanup%s()
}
`

const otherGoStub = `// +build !windows

package <INSERT PACKAGE NAME HERE>

// Platforms other than Windows do not need the shim library.

func loaderInit() {
}

func loaderDestroy() {
}
`

func writeWindowsGoStub(moduleName string, output io.Writer) error {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf(windowsGoStub, strings.ToLower(moduleName), moduleName, moduleName, moduleName, moduleName))

	_, err := output.Write(b.Bytes())
	return err
}

func writeOtherGoStub(output io.Writer) error {
	_, err := output.Write([]byte(otherGoStub))
	return err
}
