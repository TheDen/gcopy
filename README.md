# gcopy


`gcopy` (global copy) is a command line tool that copies data to your clipboard on MacOS. It's similar to `pbcopy` but with some key differences

* Works with images, so copied images can be pasted
* Accepts arbitrary `STDIN` via a pipe, or via filename passed as an argument
* Written in Go, deployed as a universal static binary

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

### Copying to clipboard via pipes

Works with text or arbitrary data

```shell
cat main.go | gcopy
```

Images can also be copied via pipes, and then pasted as images


```shell
cat image.png | gcopy
```

### Passing in files to copy to the clipboard

```shell
gcopy main.go
```

Similarly for images

```shell
gcopy image.png
```
