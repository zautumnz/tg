`tg` (Terminal in Go) is a Golang terminal emulator for MacOS.

Fork of [this fork of Darktile](https://github.com/HuckRidgeSW/darktile) which
is a fork of [Aminal](https://github.com/zautumnz/tg/tree/legacy-aminal).

## TODO:

* Fix the bug in pasteboard in clipboard (numbers are not booleans)
* Remove as much other stuff as possible, only care about this being fast and working on Mac and showing unicode
* Better default font with cjk support
* Get rid of config parsing, just need to modify and rebuild
* Get rid of transparency
* Get rid of themes
* Just one cursor style

## Key Bindings

| Action                      | Binding |
|-----------------------------|---------|
| Copy               | `ctrl + shift + C`
| Paste              | `ctrl + shift + V`
| Decrease font size | `ctrl + -`
| Increase font size | `ctrl + =`
| Open URL           | `ctrl + click`
