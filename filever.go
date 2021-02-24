package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gonutz/w32/v2"
)

func usage() {
	fmt.Print(`usage: filever [format] file

  filever prints the version information of the given executable file to
  standard ouput.
  An executable file has a major, minor, patch and build version.
  If no version info is set in the executable file, filever outputs nothing and
  returns error value 3.
  filver returns 0 on success. The output is only valid in that case.

  format  determines the output format, this is a dot-separated list of version
          names "major", "minor", "patch" and "build".
          example: "major.minor" on a file with version 1.2.3.4 outputs "1.2"
          the default value is "major.minor.patch.build"
  file    the executable file whose version is read
`)
}

const (
	errArguments     = 2
	errNoVersionInfo = 3
)

func main() {
	argCount := len(os.Args) - 1
	if !(1 <= argCount && argCount <= 2) ||
		(argCount == 1 && isHelpFlag(os.Args[1])) {
		usage()
		return
	}

	var format, exePath string
	if argCount == 1 {
		format = "major.minor.patch.build"
		exePath = os.Args[1]
	} else {
		format = strings.ToLower(os.Args[1])
		exePath = os.Args[2]
	}

	if _, err := os.Lstat(exePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(errArguments)
	}

	// get the version info from the exe file
	size := w32.GetFileVersionInfoSize(exePath)
	if size <= 0 {
		fmt.Fprintf(
			os.Stderr,
			"'%s' does not contain version info (GetFileVersionInfoSize returned %d)\n",
			exePath,
			size,
		)
		os.Exit(errNoVersionInfo)
	}
	info := make([]byte, size)
	ok := w32.GetFileVersionInfo(exePath, info)
	if !ok {
		fmt.Fprintf(
			os.Stderr,
			"'%s' does not contain version info (GetFileVersionInfo returned 0)\n",
			exePath,
		)
		os.Exit(errNoVersionInfo)
	}
	fixed, ok := w32.VerQueryValueRoot(info)
	if !ok {
		fmt.Fprintf(
			os.Stderr,
			"'%s' does not contain version info (VerQueryValueRoot returned 0)\n",
			exePath,
		)
		os.Exit(errNoVersionInfo)
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
			fmt.Fprintf(
				os.Stderr,
				"'%s' is not a valid version part\n",
				part,
			)
			os.Exit(errArguments)
		}
	}

	fmt.Print(strings.Join(output, "."))
}

func isHelpFlag(s string) bool {
	s = strings.ToLower(s)
	return s == "-h" || s == "-help" || s == "--help" || s == "/h" || s == "/?"
}
