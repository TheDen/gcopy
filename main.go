package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/h2non/filetype"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getAbsfilename(fileName string) string {
	absfileName, err := filepath.Abs(fileName)
	check(err)
	return absfileName
}

func fileCheck(filePath string) {
	fileInfo, err := os.Stat(filePath)
	check(err)
	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "error: %s is a directory\n", filePath)
		os.Exit(1)
	}
}

func getfileType(absfileName string) string {
	buf, _ := ioutil.ReadFile(absfileName)
	kind, _ := filetype.Match(buf)
	return kind.Extension
}

func getfileTypeStdin(stdin []byte) string {
	kind, _ := filetype.Match(stdin)
	return kind.Extension
}

func setfileClass(fileType string) string {
	var fileClass = ""
	switch {
	case fileType == "png":
		fileClass = "«class PNGf»"
	case fileType == "jpg":
		fileClass = "JPEG picture"
	case fileType == "gif":
		fileClass = "GIF picture"
	case fileType == "bmp":
		fileClass = "«class BMPf»"
	}
	return fileClass
}

func parseFile(absfileName string, fileType string) string {
	fileClass := setfileClass(fileType)
	var command = ""
	if len(fileClass) > 0 && (len(absfileName) > 0) {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\") as %s)", absfileName, fileClass)
	} else {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\"))", absfileName)
	}
	return command
}

func runCommand(command string) {
	cmd := exec.Command("osascript", "-e", command)
	_, err := cmd.CombinedOutput()
	check(err)
}

func writTtempFile(stdin []byte, fileType string) *os.File {
	f, err := os.CreateTemp("", fmt.Sprintf("tmpfile-*.%s", fileType))
	check(err)
	_, err = f.Write(stdin)
	check(err)
	return f
}

func main() {
	stat, _ := os.Stdin.Stat()
	// Exit if nothing is passed in
	if len(os.Args) <= 1 && ((stat.Mode() & os.ModeCharDevice) != 0) {
		os.Exit(0)
	}

	var command = ""
	var fileType = ""

	// stdin takes precedence over cli argument
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := ioutil.ReadAll(os.Stdin)
		check(err)
		fileType = getfileTypeStdin(stdin)
		tempfile := writTtempFile(stdin, fileType)
		defer tempfile.Close()
		defer os.Remove(tempfile.Name())
		command = parseFile(tempfile.Name(), fileType)
		runCommand(command)
		os.Exit(0)
	}

	// Check for argument, expecting a file
	if len(os.Args) > 1 {
		fileName := os.Args[1]
		absfileName := getAbsfilename(fileName)
		fileCheck(absfileName)
		fileType = getfileType(absfileName)
		command = parseFile(absfileName, fileType)
		runCommand(command)
		os.Exit(0)
	}
}
