package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gonutz/w32"
)

func usage() {
	fmt.Print(`usage: filever [format] file

  filever prints the version information of the given executable file. An
  executable file has a major, minor, patch and build version.
  If no version info is set in the executable file, filever outputs nothing and
  returns successfully.

  format  determines the output format, this is a dot-separated list of version
          names "major", "minor", "patch" and "build".
          example: "major.minor" on a file with version 1.2.3.4 outputs "1.2"
          the default value is "major.minor.patch.build"
  file    the executable file whose version is read
`)
}

func main() {
	argCount := len(os.Args) - 1
	if !(1 <= argCount && argCount <= 2) {
		usage()
		os.Exit(2)
	}

	var format, exePath string
	if argCount == 1 {
		format = "major.minor.patch.build"
		exePath = os.Args[1]
	} else {
		format = strings.ToLower(os.Args[1])
		exePath = os.Args[2]
	}

	if _, err := os.Stat(exePath); err != nil {
		fmt.Println("error: executable not found")
		os.Exit(2)
	}

	// get the version info from the exe file
	size := w32.GetFileVersionInfoSize(exePath)
	if size <= 0 {
		return // no version info in file, print nothing, return OK
	}
	info := make([]byte, size)
	ok := w32.GetFileVersionInfo(exePath, info)
	if !ok {
		usage()
		os.Exit(1)
	}
	fixed, ok := w32.VerQueryValueRoot(info)
	if !ok {
		usage()
		os.Exit(1)
	}
	version := fixed.FileVersion()
	major := int(version & 0xFFFF000000000000 >> 48)
	minor := int(version & 0x0000FFFF00000000 >> 32)
	patch := int(version & 0x00000000FFFF0000 >> 16)
	build := int(version & 0x000000000000FFFF >> 0)

	// parse the output format and print the version
	parts := strings.Split(format, ".")
	var output []string
	for _, part := range parts {
		switch part {
		case "major":
			output = append(output, strconv.Itoa(major))
		case "minor":
			output = append(output, strconv.Itoa(minor))
		case "patch":
			output = append(output, strconv.Itoa(patch))
		case "build":
			output = append(output, strconv.Itoa(build))
		default:
			usage()
			os.Exit(2)
		}
	}

	fmt.Print(strings.Join(output, "."))
}
