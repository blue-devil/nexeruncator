<div align="center"><img src="https://user-images.githubusercontent.com/10853207/209234630-b29fbaaa-536b-4899-8eda-3a42a2d73023.png" width=300></div>

<p align="center"></p>

<div align="center">

![License](https://img.shields.io/badge/license-GPLv3-green)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/blue-devil/nexeruncator)
[![Go Report Card](https://goreportcard.com/badge/github.com/blue-devil/nexeruncator)](https://goreportcard.com/report/github.com/blue-devil/nexeruncator)

</div>

# neXeruncator

neXeruncator; extracts or inserts javascript source file from or into nexe-compiled binaries.

* The word `neXeruncator` is derived from `nexe` and `aberuncate`. Nothing special!

## Build

neXeruncator project created using Go version 1.19.
Building project is as simple as running `go build`

```bash
go build main.go
```

## Install

neXeruncator can work as a standalone executable, you do not need to install it.
Apply the build step.

## Usage

Running `nexeruncator` prints help. There are two commands `extract` and `insert`.

```txt
          __    __                              _
 _ __   __\ \  / /__ _ __ _   _ _ __   ___ __ _| |_ ___  _ __
| '_ \ / _ \ \/ / _ | '__| | | | '_ \ / __/ _' | __/ _ \| '__|
| | | |  __/>  <  __| |  | |_| | | | | (_| (_| | || (_) | |
|_| |_|\___/ /\ \___|_|   \__,_|_| |_|\___\__,_| \_\___/|_|
          /_/  \_\                  Blue DeviL \_____/  SCT

NAME:
   nexeruncator - aberuncator of nexe-compiled binaries

USAGE:
   nexeruncator [global options] command [command options] [arguments...]

VERSION:
   0.3.0

COMMANDS:
   help     displays help messages.
   extract  extracts javascript source
   insert   inserts javascript source

GLOBAL OPTIONS:
   --version, -v  print the version (default: false)
```

### extract

`extract` command simply extracts embedded javascript file from nexe-compiled
binary.

```txt
NAME:
   nexeruncator extract - extracts javascript source

USAGE:
   nexeruncator extract [command options]

DESCRIPTION:
   Extracts javascript source file from nexe-compiled binary

OPTIONS:
   --file value  nexe-compiled binary
   --help, -h    show help (default: false)
```

Below, we can see a sample command line for extracting:

```bash
nexeruncator extract --file nexefile.exe
```

### insert

`insert` command inserts a new javascript source into nexe-compiled binary and
created a new file named: `patched.exe`. `insert` command takes 2 arguments:
nexe-compiled binary and new javascript:

```bash
NAME:
   nexeruncator insert - inserts javascript source

USAGE:
   nexeruncator insert [command options]

DESCRIPTION:
   Inserts javascript source file into nexe-compiled binary

OPTIONS:
   --dest value  nexe-compiled binary
   --src value   javascript source
   --help, -h    show help (default: false)
```

Below, we can see a sample command line for inserting javascript source into
nexe-compiled binary:

```bash
nexeruncator insert --dest nexefile.exe --src myjsfile.js
```

## Notes

From a reverse engineer's point of view, the nexe-compiled application consists of 4 main parts:

* NodeJS part, which contains NodeJS runtime and big in size, ~55mb
* Initializer javascript code.
* User's javacript code.
* 32 byte footer part.

Initializer javascript part contains name of user's javascript source filename. Let's say we have a javascript file "awesome.js". That means "awesome.js" string occurs two times in initializer part.

This is how the nexe-compiled binary ends:

```txt
Offset(h) 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F

035E3820        3C 6E 65 78 65 7E 7E 73 65 6E 74 69 6E 65    <nexe~~sentine
035E3830  6C 3E 00 00 00 00 00 FE CE 40 00 00 00 00 00 00  l>.....þÎ@......
035E3840  43 40                                            C@
```

The above sniplet is 32 bytes. The first 16 byte is the string literal `<nexe~~sentinel>`. The second 16 byte is actually comprise of 2 8 byte parts. And they are in little endian. Let's check [Nexe][gh-nexe]'s sources: `nexe/src
/
compiler.ts`

```ts
<snip>
    const code = this.code(),
        codeSize = Buffer.byteLength(code),
        lengths = Buffer.from(Array(16))

    lengths.writeDoubleLE(codeSize, 0)
    lengths.writeDoubleLE(this.bundle.size, 8)
    return new (MultiStream as any)([
        binary,
        toStream(code),
        this.bundle.toStream(),
        toStream(Buffer.concat([Buffer.from('<nexe~~sentinel>'), lengths])),
    ])
  }
}
```

The first 8 byte is `codeSize` which is what I called jsInit. It is a kind og initializer javascript code just before user's javascript code and finished with double semi-coloumns. And the second 8 byte is size of user's appended javascript source file.

## Todo

* [ ] Implement a info function
  * [ ] Collect info about nexe-compiled binary(js source name, file size)
  * [ ] Print this info
* [ ] Add to pkg.go.dev as a module
* [ ] Finish documentation
* [ ] Implement tests
* [ ] Add pre-built binaries
* [ ] Add testdata

## Resources

* [Github - Nexe][gh-nexe]

## License

This project is under GPLv3 license.

[gh-nexe]: https://github.com/nexe/nexe/
