# barnard - Cantudo fork

It adds suport for immediate conenction to a channel:
```
barnard -server mumble.FFF.org:64738 -username prueba -channel "channel"
```
support for connection to different audio devices
```
barnard -server mumble.FFF.org:64738 -username prueba -channel "channel" -outputdevice 4 -inputdevice 5
```
you can get a list of the audio devices with
```
barnard -list_devices
```
it also adds a flag to start transmission rightaway without any keypresses
```
barnard -server mumble.FFF.org:64738 -username prueba -channel "channel" -outputdevice 4 -inputdevice 5 -inmediatestart
```


barnard is a terminal-based client for the [Mumble](https://mumble.info) voice
chat software.

![Screenshot](https://i.imgur.com/B8ldT5k.png)

## Installation

Requirements:

1. [Go](https://golang.org/)
2. [Git](https://git-scm.com/)
3. [Opus](https://opus-codec.org/) development headers
4. [OpenAL](http://kcat.strangesoft.net/openal.html) development headers

To fetch and build:

    go get -u layeh.com/barnard

After running the command above, `barnard` will be compiled as `$(go env GOPATH)/bin/barnard`.

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
