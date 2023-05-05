package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	_ "embed"

	"github.com/akamensky/argparse"
	"github.com/h2non/filetype"
)

const (
	pngFileType  = "png"
	jpgFileType  = "jpg"
	gifFileType  = "gif"
	bmpFileType  = "bmp"
	tiffFileType = "tif"
)

//go:generate bash get_version.sh
//go:embed version.txt
var version string

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getAbsfilename(fileName string) string {
	absfileName, err := filepath.Abs(fileName)
	exitOnError(err)
	return absfileName
}

func validateFile(filePath string) {
	_, err := os.Stat(filePath)
	exitOnError(err)
}

func cleanup(tempfile *os.File) int {
	tempfile.Close()
	os.Remove(tempfile.Name())
	return 0
}

func getFileClass(fileContent []byte) (string, string) {
	kind, _ := filetype.Match(fileContent)
	fileExtension := kind.Extension
	isImage := filetype.IsImage(fileContent)

	var fileClass = ""
	if isImage {
		switch strings.ToLower(fileExtension) {
		case pngFileType:
			fileClass = "«class PNGf»"
		case jpgFileType:
			fileClass = "JPEG picture"
		case gifFileType:
			fileClass = "GIF picture"
		case bmpFileType:
			// Extra space here is intentional
			fileClass = "«class BMP »"
		case tiffFileType:
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

func createCommand(absfileName string, fileClass string, rawData bool) (command string) {
	if len(fileClass) > 0 {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\") as %s)", absfileName, fileClass)
	} else if rawData {
		command = fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\"))", absfileName)
	} else {
		command = fmt.Sprintf("tell application \"Finder\" to set the clipboard to (POSIX file \"%s\")", absfileName)
	}
	return
}

func runCommand(command string) {
	defer func() {
		if r := recover(); r != nil {
			// Unable to copy the data to clipboard, do nothing

		}
	}()
	cmd := exec.Command("osascript", "-e", command)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		err = fmt.Errorf("%w; %s", err, string(output))
	}
	exitOnError(err)
}

func writeTempFile(stdin []byte, fileExtension string) *os.File {
	f, err := os.CreateTemp("", fmt.Sprintf("tmpfile-*.%s", fileExtension))
	exitOnError(err)
	_, err = f.Write(stdin)
	exitOnError(err)
	return f
}

func main() {
	parser := argparse.NewParser("gcopy [file] [STDIN]", "gcopy: copy content to the clipboard")
	printVersion := parser.Flag("v", "version", &argparse.Options{Help: "Current version"})
	pathName := parser.Flag("p", "path", &argparse.Options{Help: "Copy (and show) the absolute path of a file or folder to the clipboard"})
	var fileName *string = parser.StringPositional(&argparse.Options{Help: "DISABLEDDESCRIPTIONWILLNOTSHOWUP"})
	parser.Parse(os.Args)

	if *printVersion {
		fmt.Print("build version: ", version)
		os.Exit(0)
	}

	var fileArg bool
	if len(*fileName) > 0 {
		fileArg = true
	}

	if *pathName && fileArg {
		absfileName := getAbsfilename(*fileName)
		validateFile(absfileName)
		command := fmt.Sprintf("set the clipboard to \"%s\"", absfileName)
		fmt.Println(absfileName)
		runCommand(command)
		os.Exit(0)
	}

	stat, err := os.Stdin.Stat()
	exitOnError(err)
	if !fileArg && ((stat.Mode() & os.ModeCharDevice) != 0) {
		os.Exit(0)
	}

	var command = ""
	// stdin takes precedence over positional argument
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		exitOnError(err)
		fileClass, fileExtension := getFileClass(stdin)
		tempfile := writeTempFile(stdin, fileExtension)
		command = createCommand(tempfile.Name(), fileClass, true)
		defer tempfile.Close()
		defer os.Remove(tempfile.Name())
		runCommand(command)
		os.Exit(cleanup(tempfile))
	}
	// Argument expects file
	if fileArg {
		absfileName := getAbsfilename(*fileName)
		validateFile(absfileName)
		command = createCommand(absfileName, "", false)
		runCommand(command)
		os.Exit(0)
	}
}
