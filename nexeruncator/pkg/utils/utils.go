// Package utils contains main logic of extracting and inserting javascript
// file from nexe-compiled binaries
package utils

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

var MagicBytes = []byte{
	0x21, 0x28, 0x66, 0x75, 0x6E, 0x63, 0x74, 0x69, 0x6F, 0x6E, 0x20, 0x28,
	0x29, 0x20, 0x7B, 0x70, 0x72, 0x6F, 0x63, 0x65, 0x73, 0x73, 0x2E, 0x5F,
	0x5F, 0x6E, 0x65, 0x78, 0x65, 0x20, 0x3D, 0x20, 0x7B, 0x22, 0x72, 0x65,
	0x73, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x73, 0x22, 0x3A, 0x7B, 0x22, 0x2E,
	0x2F,
}

// Checks OS and slears command prompt/terminal
func ClearConsole() {
	ros := runtime.GOOS
	switch ros {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "darwin":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Printf("[-] OS %s not defined", ros)
	}
}

// CheckNoColor returns true if NO_COLOR environmental variable is set
func CheckNoColor() bool {
	_, exist := os.LookupEnv("NO_COLOR")
	return exist
}

// Prints nexeruncator banner to command prompt/terminal
func PrintBanner() {
	banner := `
          __    __                              _             
 _ __   __\ \  / /__ _ __ _   _ _ __   ___ __ _| |_ ___  _ __ 
| '_ \ / _ \ \/ / _ | '__| | | | '_ \ / __/ _' | __/ _ \| '__|
| | | |  __/>  <  __| |  | |_| | | | | (_| (_| | || (_) | |   
|_| |_|\___/ /\ \___|_|   \__,_|_| |_|\___\__,_| \_\___/|_|   
          /_/  \_\                  Blue DeviL \_____/  SCT
`

	ClearConsole()
	if CheckNoColor() {
		fmt.Println(banner)
	} else {
		colorBanner := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, banner)
		fmt.Println(colorBanner)
	}
}

// Returns javascript source file's name
func getFilename(slice []byte) string {
	// search for MagicBytes in slice
	index := bytes.Index(slice, MagicBytes)
	if index == -1 {
		return ""
	}
	// MagicBytes found, create new slice starting at index + len(MagicBytes)
	filenameSlice := slice[index+len(MagicBytes):]
	// search for double quote in filenameSlice
	index = bytes.IndexByte(filenameSlice, '"')
	if index == -1 {
		return ""
	}
	// double quote found, create new slice from start to index
	filenameSlice = filenameSlice[:index]
	// convert slice to string and return
	return string(filenameSlice)
}

// Searches for javascript source inside nexe-compiled binary and if found
// saves the file as its own name
func ExtractJS(filename string) error {
	// read file into nexeFile slice
	nexeFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[-] Error reading file: %s\n", err)
		os.Exit(1)
	}
	// search for MagicBytes in nexeFile
	index := bytes.Index(nexeFile, MagicBytes)
	if index == -1 {
		fmt.Println("[-] Are you sure if this is a nexe compiled binary?")
		os.Exit(1)
	}
	// MagicBytes found, store index in nexeStart
	nexeStart := index
	// create jsWhole slice from nexeStart to end of nexeFile
	jsWhole := nexeFile[nexeStart:]
	// get filename from jsWhole slice
	jsFilename := getFilename(jsWhole)
	if jsFilename == "" {
		fmt.Println("[-] Cannot get filename...")
		os.Exit(1)
	}
	// search for ";;" in jsWhole slice
	index = bytes.Index(jsWhole, []byte(";;"))
	if index == -1 {
		fmt.Println("[-] Error: cannot find ;;")
		os.Exit(1)
	}

	// ";;" found, create script2 slice from index + 2 to end of jsWhole
	script2 := jsWhole[index+2:]
	// search for <nexe~~sentinel> in script2 slice
	index = bytes.Index(script2, []byte("<nexe~~sentinel>"))
	if index == -1 {
		fmt.Println("[-] Cannot find <nexe~~sentinel>...")
		os.Exit(1)
	}
	// <nexe~~sentinel> found, create finalScript slice from start of script2 to index
	finalScript := script2[:index]
	// save finalScript slice as binary file with jsFilename
	err = os.WriteFile(jsFilename, finalScript, 0644)
	if err != nil {
		fmt.Printf("[-] Error saving file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("[+] Extracted script saved as %s\n", jsFilename)

	return nil
}

// Sets javascript sources's name and its size inside nexe-compiled binary.
func setFilenameAndSize(slice []byte, fileName string, size uint64) ([]byte, uint64) {
	magicBytes := []byte{0x3A, 0x7B, 0x22, 0x2E, 0x2F}
	// search for magicBytes in slice
	index := bytes.Index(slice, magicBytes)
	if index == -1 {
		// magicBytes not found
		return nil, 0
	}
	// magicBytes found, create retVal slice starting at index + len(magicBytes)
	retVal := make([]byte, 16000)
	copy(retVal, slice[:index+len(magicBytes)])
	// convert fileName to slice and append to retVal
	retVal = append(retVal, []byte(fileName)...)
	// append [0x22, 0x3A, 0x5B, 0x30, 0x2C] to retVal
	retVal = append(retVal, []byte{0x22, 0x3A, 0x5B, 0x30, 0x2C}...)
	// convert size to string and then to slice, append to retVal
	sizeStr := strconv.FormatUint(uint64(size), 10)
	retVal = append(retVal, []byte(sizeStr)...)
	// search for [0x5D, 0x7D, 0x7D, 0x3B, 0x0A, 0x7D, 0x29] in slice
	index = bytes.Index(slice, []byte{0x5D, 0x7D, 0x7D, 0x3B, 0x0A, 0x7D, 0x29})
	if index == -1 {
		// [0x5D, 0x7D, 0x7D, 0x3B, 0x0A, 0x7D, 0x29] not found
		return nil, 0
	}

	// search 0x65, 0x78, 0x65, 0x63, 0x50, 0x61, 0x74, 0x68, 0x29, 0x2C, 0x22, 0x2E, 0x2F
	index2 := bytes.Index(slice, []byte{0x65, 0x78, 0x65, 0x63, 0x50, 0x61, 0x74, 0x68, 0x29, 0x2C, 0x22, 0x2E, 0x2F})
	if index2 == -1 {
		// [0x65, 0x78, 0x65, 0x63, 0x50, 0x61, 0x74, 0x68, 0x29, 0x2C, 0x22, 0x2E, 0x2F] not found
		return nil, 0
	}

	// [0x5D, 0x7D, 0x7D, 0x3B, 0x0A, 0x7D, 0x29] found, append from index to end of slice to retVal
	temp := slice[index:(index2 + 13)]
	retVal = append(retVal, temp...)
	// convert fileName to slice and append to retVal
	retVal = append(retVal, []byte(fileName)...)

	//  search 0x22, 0x29, 0x0A, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x70
	index2 = bytes.Index(slice, []byte{0x22, 0x29, 0x0A, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x70})
	if index2 == -1 {
		// [0x65, 0x78, 0x65, 0x63, 0x50, 0x61, 0x74, 0x68, 0x29, 0x2C, 0x22, 0x2E, 0x2F] not found
		return nil, 0
	}
	retVal = append(retVal, slice[index2:]...)

	// return retVal and its size
	return retVal, uint64(len(retVal))
}

// Inserts and overwrites a new javascript file into a nexe-compiled binary
func InsertJS(nexeFile string, scriptFile string) error {
	// read script and nexe files as slices of bytes
	jsscript, err := os.ReadFile(scriptFile)
	if err != nil {
		fmt.Printf("[-] Error reading script file: %s\n", err)
		os.Exit(1)
	}
	wholeProgram, err := os.ReadFile(nexeFile)
	if err != nil {
		fmt.Printf("[-] Error reading nexe file: %s\n", err)
		os.Exit(1)
	}
	// store script file size and name
	lenScript := uint64(len(jsscript))
	nameScript := scriptFile

	// create engine slice from start of wholeProgram to magicBytes
	index := bytes.Index(wholeProgram, MagicBytes)
	if index == -1 {
		fmt.Println("[-] Are you sure if this is a nexe compiled binary?")
		os.Exit(1)
	}
	engine := wholeProgram[:index]

	// create jsWhole slice from index to end of wholeProgram
	jsWhole := wholeProgram[index:]
	// search for ";;" in jsWhole
	index = bytes.Index(jsWhole, []byte(";;"))
	if index == -1 {
		fmt.Println("[-] Error: cannot find ;;")
		os.Exit(1)
	}
	// ";;" found, create jsInit slice from start of jsWhole to index + 2
	jsInit := jsWhole[:index+2]

	//
	// NOTE(bluedevil): Gather info about files

	// store size of jsInit
	// oldCodeSize := uint64(len(jsInit))
	// fmt.Println(oldCodeSize)

	// call setFilenameAndSize with jsInit, nameScript, and lenScript
	jsInitFin, codeSize := setFilenameAndSize(jsInit, nameScript, lenScript)

	// convert lenScript and codeSize to float64 and make 8-byte slices
	lenScriptF64 := math.Float64bits(float64(lenScript))
	lenScriptBytes := []byte{
		byte(lenScriptF64 >> 56),
		byte(lenScriptF64 >> 48),
		byte(lenScriptF64 >> 40),
		byte(lenScriptF64 >> 32),
		byte(lenScriptF64 >> 24),
		byte(lenScriptF64 >> 16),
		byte(lenScriptF64 >> 8),
		byte(lenScriptF64),
	}
	// reverse and append lenScriptBytes and codeSizeBytes to create footer
	codeSizeF64 := math.Float64bits(float64(codeSize))
	codeSizeBytes := []byte{
		byte(codeSizeF64 >> 56),
		byte(codeSizeF64 >> 48),
		byte(codeSizeF64 >> 40),
		byte(codeSizeF64 >> 32),
		byte(codeSizeF64 >> 24),
		byte(codeSizeF64 >> 16),
		byte(codeSizeF64 >> 8),
		byte(codeSizeF64),
	}
	for i, j := 0, len(lenScriptBytes)-1; i < j; i, j = i+1, j-1 {
		lenScriptBytes[i], lenScriptBytes[j] = lenScriptBytes[j], lenScriptBytes[i]
	}
	for i, j := 0, len(codeSizeBytes)-1; i < j; i, j = i+1, j-1 {
		codeSizeBytes[i], codeSizeBytes[j] = codeSizeBytes[j], codeSizeBytes[i]
	}
	footer := []byte("<nexe~~sentinel>")
	footer = append(footer, codeSizeBytes...)
	footer = append(footer, lenScriptBytes...)

	// create patched slice by concatenating engine, jsInitFin, jsscript, and footer
	patched := append(engine, jsInitFin...)
	patched = append(patched, jsscript...)
	patched = append(patched, footer...)

	// write patched slice to patched.exe file
	err = os.WriteFile("patched.exe", patched, 0644)
	if err != nil {
		fmt.Printf("[-] Error writing patched file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("[+] File %s inserted into %s\n", scriptFile, nexeFile)

	return nil
}
