# go-fu

Search [commandlinefu.com](http://www.commandlinefu.com/) from CLI.

## Installation

    go get github.com/t9md/go-fu

## Usage

    go-fu COMMAND [PAGE]
    
      COMMAND: browse, using WORD, by USER, matching WORD
      PAGE: 1-999 (defaut: 1)

## Example

    go-fu browse
    go-fu using grep
    go-fu by USER
    go-fu matching find 2

## Abbreviation

You can abbreviate COMMAND like follwing.

    go-fu br
    go-fu u grep
    go-fu by USER
    go-fu m find 2
