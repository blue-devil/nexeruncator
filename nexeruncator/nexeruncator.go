// Package nexeruncator provides two functionalities. nexeruncator can
// extract embedded javascript source from nexe-compiled binaries.
// nexeruncator can insert a new javascript source to a nexe-compiled
// binary
package nexeruncator

import (
	"log"

	"github.com/blue-devil/nexeruncator/nexeruncator/pkg/utils"
)

// ExtractJavascript function extracts javascript source file from nexeFile
func ExtractJavascript(nexeFile string) {
	var err error
	if err = utils.ExtractJS(nexeFile); err != nil {
		log.Printf("[-] Failed to extract javasript source file from %s, %v", nexeFile, err)
	}
}

// InsertJavascript function inserts jsFile into nexefile. Both nexeFile
// and jsFile is the name of files. Function did not overwrites the original
// executable file, creates a new file with name `pathced.exe`
func InsertJavascript(nexeFile, jsFile string) {
	var err error
	if err = utils.InsertJS(nexeFile, jsFile); err != nil {
		log.Printf("[-] Failed to insert javascript source file %s to %s, error: %v", jsFile, nexeFile, err)
	}
}
