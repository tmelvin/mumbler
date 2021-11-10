package main

import (
	"crypto/tls"

	"github.com/acien101/barnard/uiterm"
	"github.com/cantudo/gumble/gumble"
	"github.com/cantudo/gumble/gumbleopenal"
)

type Barnard struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	Channel   string
	TLSConfig tls.Config

	Stream *gumbleopenal.Stream

	Ui            *uiterm.Ui
	UiOutput      uiterm.Textview
	UiInput       uiterm.Textbox
	UiStatus      uiterm.Label
	UiTree        uiterm.Tree
	UiInputStatus uiterm.Label
}
