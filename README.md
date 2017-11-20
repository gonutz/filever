filever
=======

```
usage: filever [format] file

  filever prints the version information of the given executable file. An
  executable file has a major, minor, patch and build version.
  If no version info is set in the executable file, filever outputs nothing and
  returns successfully.

  format  determines the output format, this is a dot-separated list of version
          names "major", "minor", "patch" and "build".
          example: "major.minor" on a file with version 1.2.3.4 outputs "1.2"
          the default value is "major.minor.patch.build"
  file    the executable file whose version is read
```

Build
-----

To build and install this tool, use:

`go get github.com/gonutz/filever`

There is also a build script `build.bat` which builds a small (stripping symbol and debug info from the result) 32 bit application instead of using the default Go behavior which builds 64 bit executables on 64 bit machines and always includes symbol and debug info.