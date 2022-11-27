package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf8"

	_ "embed"

	"github.com/akamensky/argparse"
	"github.com/h2non/filetype"
)

//go:generate bash get_version.sh
//go:embed version.txt
var version string

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func checkErrExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getAbsfilename(fileName string) string {
	absfileName, err := filepath.Abs(fileName)
	checkErrExit(err)
	return absfileName
}

func fileCheck(filePath string) {
	fileInfo, err := os.Stat(filePath)
	checkErrExit(err)
	if fileInfo.IsDir() {
		checkErrExit(err)
	}
}

func cleanup(tempfile *os.File) int {
	tempfile.Close()
	os.Remove(tempfile.Name())
	return 0
}

func getfileClass(fileContent []byte) (string, string) {
	kind, _ := filetype.Match(fileContent)
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
			// Extra space here is intentional
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

func createCommand(absfileName string, fileClass string, rawData bool) string {
	if len(fileClass) > 0 {
		return fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\") as %s)", absfileName, fileClass)
	}
	if rawData {
		return fmt.Sprintf("set the clipboard to (read (POSIX file \"%s\"))", absfileName)
	} else {
		return fmt.Sprintf("tell application \"Finder\" to set the clipboard to (POSIX file \"%s\")", absfileName)
	}
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
	checkPanic(err)
}

func writeTempFile(stdin []byte, fileExtension string) *os.File {
	f, err := os.CreateTemp("", fmt.Sprintf("tmpfile-*.%s", fileExtension))
	checkPanic(err)
	_, err = f.Write(stdin)
	checkPanic(err)
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
		fileCheck(absfileName)
		command := fmt.Sprintf("set the clipboard to \"%s\"", absfileName)
		fmt.Println(absfileName)
		runCommand(command)
		os.Exit(0)
	}

	stat, _ := os.Stdin.Stat()
	if !fileArg && ((stat.Mode() & os.ModeCharDevice) != 0) {
		os.Exit(0)
	}

	var command = ""
	// stdin takes precedence over positional argument
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		checkErrExit(err)
		fileClass, fileExtension := getfileClass(stdin)
		tempfile := writeTempFile(stdin, fileExtension)
		command = createCommand(tempfile.Name(), fileClass, true)
		defer tempfile.Close()
		defer os.Remove(tempfile.Name())
		runCommand(command)
		os.Exit(cleanup(tempfile))
	}
	// Check for argument, expecting a file
	if fileArg {
		absfileName := getAbsfilename(*fileName)
		fileCheck(absfileName)
		command = createCommand(absfileName, "", false)
		runCommand(command)
		os.Exit(0)
	}
}
