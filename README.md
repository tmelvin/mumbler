# Mumbler - Barnard fork

You'll very likely want to install Barnards fork of Cantudo:

https://github.com/layeh/barnard

It adds suport for immediate conenction to a channel:
```
mumbler -server mumble.FFF.org:64738 -username prueba -channel "channel"
```
support for connection to different audio devices
```
mumbler -server mumble.FFF.org:64738 -username prueba -channel "channel" -outputdevice 4 -inputdevice 5
```
you can get a list of the audio devices with
```
mumbler -list_devices
```
it also adds a flag to start transmission rightaway without any keypresses
```
mumbler -server mumble.FFF.org:64738 -username prueba -channel "channel" -outputdevice 4 -inputdevice 5 -inmediatestart
```

## Installation

Requirements:

1. [Go](https://golang.org/)
2. [Git](https://git-scm.com/)
3. [Opus](https://opus-codec.org/) development headers
4. [OpenAL](http://kcat.strangesoft.net/openal.html) development headers

To fetch and build:

    go get -u https://github.com/tmelvin/mumbler@latest

After running the command above, `mumbler` will be compiled as `$(go env GOPATH)/bin/mumbler`.

## Manual

### Key bindings

- <kbd>F1</kbd>: toggle voice transmission
- <kbd>Ctrl+L</kbd>: clear chat log
- <kbd>Tab</kbd>: toggle focus between chat and user tree
- <kbd>Page Up</kbd>: scroll chat up
- <kbd>Page Down</kbd>: scroll chat down
- <kbd>Home</kbd>: scroll chat to the top
- <kbd>End</kbd>: scroll chat to the bottom
- <kbd>F10</kbd>: quit

## License

GPLv2

## Author

Tim Cooper (<tim.cooper@layeh.com>)
