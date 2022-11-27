# gcopy


`gcopy` (global copy) is a command line tool that copies data to your clipboard on MacOS. It's similar to `pbcopy` but with some key differences

* Works with images, so copied images can be pasted
* Accepts arbitrary `STDIN` via a pipe, or via filename passed as an argument
* Written in Go, deployed as a universal static binary
* Has an optional flag to copy the absolute pathname of a file or folder to the clipboard

![gcopy](./gcopy.gif)


# Install

Via `go install`

```shell
go install github.com/TheDen/gcopy@latest
```

Or if you want to download the latest binary for both intel and Apple silicon:

```shell
curl -sL -o gcopy 'https://github.com/TheDen/gcopy/releases/latest/download/gcopy' && chmod +x gcopy && mv gcopy /usr/local/bin/
```

For specific architecture:

* `https://github.com/TheDen/gcopy/releases/latest/download/gcopy-darwin-arm64` (Apple silicon)
* `https://github.com/TheDen/gcopy/releases/latest/download/gcopy-darwin-amd64` (intel)

## Usage

```bash
usage: gcopy [file] [STDIN] [-h|--help] [-v|--version] [-p|--path]

                            gcopy: copy content to the clipboard

Arguments:

  -h  --help     Print help information
  -v  --version  Current version
  -p  --path     Copy (and show) the absolute pathname of a file or folder to
                 the clipboard
```


### Copying to clipboard via pipes

Works with text or arbitrary data

```shell
cat main.go | gcopy
# or
gcopy < main.go
```

Images can also be copied via pipes, and then pasted as images to GUI applications


```shell
cat image.png | gcopy
# or 
gcopy < image.png
```

### Passing in files to copy to the clipboard

```shell
gcopy main.go
```

Similarly for images

```shell
gcopy image.pnggi
```

### Getting the absolute path of a file or folder

```shell
gcopy -p .bashrc
/Users/den/.bashrc
```