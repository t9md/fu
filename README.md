# go-cmdline-fu
Now development
go version of cmdline-fu

## Install
    go get github.com/t9md/go-cmdline-fu

## Update
    go get -u github.com/t9md/go-cmdline-fu

## Usage

    go-cmdline-fu COMMAND [PAGE]
    
      COMMAND: browse, using WORD, by USER, matching WORD
      PAGE: 1-999 (defaut: 1)

## Example

    go-cmdline-fu browse
    go-cmdline-fu using grep
    go-cmdline-fu by USER
    go-cmdline-fu matching find 2

## Abbreviation
    you can abbreviate COMMAND

    go-cmdline-fu br
    go-cmdline-fu u grep
    go-cmdline-fu by USER
    go-cmdline-fu m find 2
