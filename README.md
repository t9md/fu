# fu

fu is command utility to search [commandlinefu.com](http://www.commandlinefu.com/) from CLI.

## Installation

    go get github.com/t9md/fu

## Usage

    fu COMMAND [PAGE]
    
      COMMAND: browse, using WORD, by USER, matching WORD
      PAGE: 1-999 (defaut: 1)

## Example

    fu browse
    fu using grep
    fu by USER
    fu matching find 2

## Abbreviation

You can abbreviate COMMAND like follwing.

    fu br
    fu u grep
    fu by USER
    fu m find 2
