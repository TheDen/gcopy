# gcopy

<p align="center">
  <img alt="copy Logo" src="https://raw.githubusercontent.com/TheDen/gcopy/main/.github/logo.png" />
</p>

<div align="center">

[![License](https://img.shields.io/github/license/TheDen/gcopy?style=flat-square)](/LICENCE)
[![GitHub Release](https://img.shields.io/github/v/release/TheDen/gcopy?style=flat-square)](https://github.com/TheDen/gcopy/releases/latest)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/TheDen/gcopy)

</div>

`gcopy` (global copy) is a command line tool that copies data to your clipboard on MacOS.

It does what `pbcopy` does but with some extra features:

* Works with images, so copied images can be pasted in other applications
* Accepts arbitrary `STDIN` via a pipe
* If a filename is passed as an argument it'll copy the data as though you hit `CMD+C` via `Finder`, allowing pasting of files (PDFs, archive files, videos etc.) to other programs
* Has an optional flag to copy the absolute pathname of a file or folder to the clipboard
* Written in Go, deployed as a multi-arch static binary
* Works as a drop-in replacement for `pbcopy` since it extends its features

![gcopy](./gcopy-usage-example.gif)

# Install

## Via `go install`

```shell
go install github.com/TheDen/gcopy@latest
```

## Brew

To install via `brew`

```shell
brew tap theden/gcopy
brew install gcopy
```

# Aliasing

`gcopy` can be used as a drop-in replacement for `pbcopy`. To do this, add the `alias` into your `.bashrc` or similar

```shell
alias pbcopy='gcopy'
```

# Usage

## Overview

```bash
usage: gcopy [file] [STDIN] [-h|--help] [-v|--version] [-p|--path]

                            gcopy: copy content to the clipboard

Arguments:

  -h  --help     Print help information
  -v  --version  Current version
  -p  --path     Copy (and show) the absolute pathname of a file or folder to
                 the clipboard
```


## Copying to clipboard via STDIN

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

## Passing in files to copy to the clipboard

```shell
gcopy main.go
```

Similarly for images

```shell
gcopy image.png
```

Or any other type of file

```shell
gcopy backups.zip
```

Copying via this method will allow you to paste non-text data in other applications


## Getting the absolute path of a file or folder

```shell
gcopy -p .bashrc
/Users/den/.bashrc
```
