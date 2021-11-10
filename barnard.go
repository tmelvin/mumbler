package main

import (
	"crypto/tls"

	"github.com/cantudo/barnard/gumble/gumble"
	"github.com/cantudo/barnard/gumble/gumbleopenal"
	"github.com/cantudo/barnard/uiterm"
)

type Barnard struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	Channel   string
	TLSConfig tls.Config

	InputDevice  string
	OutputDevice string

	Stream *gumbleopenal.Stream

	Ui            *uiterm.Ui
	UiOutput      uiterm.Textview
	UiInput       uiterm.Textbox
	UiStatus      uiterm.Label
	UiTree        uiterm.Tree
	UiInputStatus uiterm.Label
}
