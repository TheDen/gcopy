package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf8"

	"github.com/h2non/filetype"
)

func check(err error) {
	if err != nil {
		panic(err)
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
		panic(err)
	}
}

func cleanup(tempfile *os.File) int {
	tempfile.Close()
	os.Remove(tempfile.Name())
	return 0
}

func getfileClass(fileContent []byte) (string, string) {
	kind, err := filetype.Match(fileContent)
	check(err)
	fileExtension := kind.Extension
	isImage := filetype.IsImage(fileContent)

	var fileClass = ""
	if isImage {
		switch {
		case fileExtension == "png":
			fileClass = "«class PNGf»"
		case fileExtension == "jpg":
			fileClass = "JPEG picture"
		case fileExtension == "gif":
			fileClass = "GIF picture"
		case fileExtension == "bmp":
			fileClass = "«class BMP »"
		case fileExtension == "tif":
			fileClass = "TIFF picture"
		}
	} else {
		// Check if file is utf8 encoded
		if utf8.Valid(fileContent) {
			fileClass = "«class utf8»"
		}
	}
	return fileClass, fileExtension
}

func createCommand(absfileName string, fileClass string) string {
	var command = ""
	if len(fileClass) > 0 {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\") as %s)", absfileName, fileClass)
	} else {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\"))", absfileName)
	}
	return command
}

func runCommand(command string) {
	defer func() {
		if r := recover(); r != nil {
			// Unable to cooy the data to clipboard, do nothing
		}
	}()
	cmd := exec.Command("osascript", "-e", command)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		err = fmt.Errorf("%w; %s", err, string(output))
	}
	check(err)
}

func writTtempFile(stdin []byte, fileExtension string) *os.File {
	f, err := os.CreateTemp("", fmt.Sprintf("tmpfile-*.%s", fileExtension))
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

	// stdin takes precedence over cli argument
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		check(err)
		fileClass, fileExtension := getfileClass(stdin)
		tempfile := writTtempFile(stdin, fileExtension)
		command = createCommand(tempfile.Name(), fileClass)
		defer tempfile.Close()
		defer os.Remove(tempfile.Name())
		runCommand(command)
		os.Exit(cleanup(tempfile))
	}
	// Check for argument, expecting a file
	if len(os.Args) > 1 {
		fileName := os.Args[1]
		absfileName := getAbsfilename(fileName)
		fileCheck(absfileName)
		fileContent, err := os.ReadFile(absfileName)
		check(err)
		fileClass, _ := getfileClass(fileContent)
		command = createCommand(absfileName, fileClass)
		runCommand(command)
		os.Exit(0)
	}
}
